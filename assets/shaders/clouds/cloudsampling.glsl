const int   CLOUD_DOMAIN = 20000;
const float BASE_DENSITY = 0.5;

vec2 boundUV(in vec3 pos, vec2 size) {
    float x = pos.x;
    float y = pos.z;
    
    while(x < 0.0) { x += size.x; }
    while(y < 0.0) { y += size.y; }

    x = mod(x, size.x) / size.x;
    y = mod(y, size.y) / size.y;

    return vec2(x, y);
}

// calculates the relative height of a sample position pos within the atmosphere defined
// by the inner and outer height of the athmosphere layers while the relative height is clamped
// between 0 and 1
float relativeHeight(in vec3 pos) {
    float len = length(pos.y - innerHeight) / (outerHeight - innerHeight);
    return clamp(len, 0.0, 1.0);
}

float cloudRemap(float h, float a, float b, float c) {
    return clampRemap(h, 0, a, 0, 1) * clampRemap(h, b, c, 1, 0);
}

float heightGradient(float cloudType, float h) {
    // calc cloud gradients
    float a = cloudRemap(h, 0.1, 0.2, 0.3);
    float b = cloudRemap(h, 0.2, 0.3, 0.5);
    float c = cloudRemap(h, 0.1, 0.7, 0.8);
    // calc weights
    vec3 weights = lerp3(h);
    float cloudGradient = a*weights.x + b*weights.y + c*weights.z;
    return cloudGradient * h;
}

vec4 sampleCloudMap(in vec3 pos) {
    vec3 p = pos / CLOUD_DOMAIN;
    return texture(cloudMapTex, p.xz);
}

float sampleCloudBase(in vec3 pos, float h) {
    // grab value from texture
    vec3 p = pos / CLOUD_DOMAIN;
    vec4 cloudBase = texture(cloudBaseTex, vec3(p.xz, h));

    // calculate low freq fbm
    float lowFreqNoise = cloudBase.r;
    float highFreqNoise = dot(cloudBase.gba, vec3(0.625, 0.25, 0.125));
    float baseDensity = clampRemap(lowFreqNoise, highFreqNoise, 1.0, 0.0, 1.0);
    baseDensity = clamp(baseDensity, 0, 1);

    // extract cloud information
    vec4  cloudInfo = sampleCloudMap(pos);
    float coverage  = cloudInfo.r;
    float cloudType = cloudInfo.b;

    // apply height gradient
    baseDensity *= heightGradient(cloudType, h);
    baseDensity *= coverage;

    return baseDensity;
}

// returns the cloud density at the specified position
float density(in vec3 pos, in float height) {
    // offset position by wind direction and strength
    vec3 newPos = pos;
    
    // get base density
    float baseDensity = sampleCloudBase(pos, height) * BASE_DENSITY;

    // perform expensive sampling if ray is within cloud
    if(baseDensity > 0.0) {
        // calculate cloud detail density and erode it from
        // the base density
        // TODO
    }
    
    return baseDensity;
}

// returns the cloud depth around a position
float getLODCloudDepth(in vec3 pos, float ds) {
    // get sample positions around the position
    vec3 p1 = pos + vec3(ds, 0, 0);
    vec3 p2 = pos - vec3(ds, 0, 0);
    vec3 p3 = pos + vec3( 0,ds, 0);
    vec3 p4 = pos - vec3( 0,ds, 0);
    vec3 p5 = pos + vec3( 0, 0,ds);
    vec3 p6 = pos - vec3( 0, 0,ds);

    // get relative height
    float h1 = relativeHeight(p1);
    float h2 = relativeHeight(p2);
    float h3 = relativeHeight(p3);
    float h4 = relativeHeight(p4);
    float h5 = relativeHeight(p5);
    float h6 = relativeHeight(p6);

    // sum depth for all sampled positions
    float sum = 0;
    sum += density(p1, h1);
    sum += density(p2, h2);
    sum += density(p3, h3);
    sum += density(p4, h4);
    sum += density(p5, h5);
    sum += density(p6, h6);
    sum /= 6;

    return sum;
}