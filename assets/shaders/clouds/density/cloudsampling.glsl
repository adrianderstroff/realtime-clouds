const int   CLOUD_DOMAIN = 4000;
uniform float uGlobalDensity = 0.2;

// calculates the relative height of a sample position pos within the atmosphere defined
// by the inner and outer height of the athmosphere layers while the relative height is clamped
// between 0 and 1
float relativeHeight(in vec3 pos) {
    return clampRemap(pos.y, innerHeight, outerHeight, 0, 1);
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
    vec3 weights = lerp3(cloudType);
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
    float baseDensity = sampleCloudBase(pos, height) * uGlobalDensity;

    // perform expensive sampling if ray is within cloud
    if(baseDensity > 0.0) {
        // calculate cloud detail density and erode it from
        // the base density
        // TODO
    }
    
    return baseDensity;
}