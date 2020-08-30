#include "../util/math.glsl"

const float CLOUD_LAYER_WIDTH = 40000;

vec3 loop(in vec3 pos, float bounds) {
    /*
    vec3 temp = mod(pos/bounds, 1);
    if(temp.x < 0) temp.x = 1 + temp.x;
    if(temp.y < 0) temp.y = 1 + temp.y;
    if(temp.z < 0) temp.z = 1 + temp.z;
    */
    return pos/bounds;
}

float density(in vec3 pos, float h, in vec3 windDir, float time, float globalCoverage, float globalDensity) {
    vec3 p = loop(pos, CLOUD_LAYER_WIDTH);

    // calculate the sample offset based on the wind direction and height. clouds move faster the higher in the cloud
    // space we sample. also because the textures we use are finite, we have to loop the sample offset to get an offset
    // that is in bounds of the cloud map texture
    vec3 off = windDir*time;
    off += windDir*h*500;
    vec3 poff = loop(pos + off, CLOUD_LAYER_WIDTH);
    // sample the texture at the offset position and then extract the different color components of the sample as the
    // cloud coverage for low and hight frequencies as well as the type of the cloud
    
    // TODO read again what low and high coverage actually means
    vec4  cloudInfo = texture(cloudMapTex, poff.xz);
    float lowCoverage  = cloudInfo.r;
    float highCoverage = cloudInfo.g;
    float cloudType    = cloudInfo.b;

    // calculate probability that clouds will form
    float cloudProbability = max(lowCoverage, saturate(globalCoverage-0.5) * 2*highCoverage);

    // calculate low freq fbm
    vec4 cloudBase = texture(cloudBaseTex, vec3(p.xz, h));
    float lowFreqNoise = cloudBase.r;
    float highFreqNoise = dot(cloudBase.gba, vec3(0.625, 0.25, 0.125));
    float baseDensity = clampRemap(lowFreqNoise, highFreqNoise-1, 1.0, 0.0, 1.0);
    baseDensity = clamp(baseDensity, 0, 1);

    // height gradient
    float heightGradient = globalDensity;

    // calculate the shape noise
    float shapeNoise = saturate(remap(baseDensity, 1 - globalCoverage*cloudProbability, 1, 0, 1)) * heightGradient;

    return shapeNoise;
}