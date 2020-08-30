#version 430
//--------------------------------------------------------------------------------------------------------------------//
// textures                                                                                                           //
//--------------------------------------------------------------------------------------------------------------------//
layout(binding = 0) uniform sampler3D cloudBaseTex;
layout(binding = 1) uniform sampler3D cloudDetailTex;
layout(binding = 2) uniform sampler2D turbulenceTex;
layout(binding = 3) uniform sampler2D cloudMapTex;

//--------------------------------------------------------------------------------------------------------------------//
// includes                                                                                                           //
//--------------------------------------------------------------------------------------------------------------------//
#include "util/camera.glsl"
#include "util/ray.glsl"
#include "util/math.glsl"
//#include "cloud/density.glsl"

//--------------------------------------------------------------------------------------------------------------------//
// uniforms                                                                                                           //
//--------------------------------------------------------------------------------------------------------------------//
// camera
uniform Camera uCamera;
// sun
uniform vec3   uSunPos           = vec3(40000, -1000, 0);
// atmosphere
uniform float  uInnerHeight      = 14000;
uniform float  uOuterHeight      = 40000;
uniform float  uExtinctionCoeff  = 1.0/26000.0;
// clouds
uniform float  uGlobalDensity    = 0.5;
uniform float  uGlobalCoverage   = 0.5;
// animation
uniform float  uTime             = 0;
uniform float  uWindSpeed        = 10;
uniform vec3   uWindDir          = vec3(1, 0, 0);
// colors
uniform vec3   uSunColor         = vec3(1, 1, 0);
uniform vec3   uAmbientColor     = vec3(1, 0, 0);
uniform vec3   uAtmosphereColor  = vec3(0.6, 0.7, 0.95);

//--------------------------------------------------------------------------------------------------------------------//
// constants                                                                                                          //
//--------------------------------------------------------------------------------------------------------------------//
const float CLOUD_LAYER_WIDTH = 150000;

//--------------------------------------------------------------------------------------------------------------------//
// input                                                                                                              //
//--------------------------------------------------------------------------------------------------------------------//
in Vertex {
    vec2 uv;
} i;

//--------------------------------------------------------------------------------------------------------------------//
// output                                                                                                             //
//--------------------------------------------------------------------------------------------------------------------//
out vec4 fragColor;

//--------------------------------------------------------------------------------------------------------------------//
// helper functions                                                                                                   //
//--------------------------------------------------------------------------------------------------------------------//
float intersectLayer(in Ray ray, in float h) {
    return (h - ray.o.y) / ray.dir.y;
}

vec3 loop(in vec3 pos, float bounds) {
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
    // sample the cloud map at the offset position and then extract the different color components of the sample as the
    // cloud coverage for low and high frequencies as well as the type of the cloud. according to the 2017 presentation
    // (and therefore according to the book The Clouds by Clausse and Facy) clouds form where air masses overlap
    
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

//--------------------------------------------------------------------------------------------------------------------//
// entry point                                                                                                        //
//--------------------------------------------------------------------------------------------------------------------//
void main() {
    // determine ray direction
    Ray ray = calcRay(i.uv, uCamera);

    // get start and end points of the horizon
    float tInner = intersectLayer(ray, uInnerHeight);
    float tOuter = intersectLayer(ray, uOuterHeight);

    // step size
    float stepSize = (tOuter - tInner) / 40.0;

    // setup ray marching variables
    float alpha         = 0.0;
    float transmittance = 1.0;
    vec3  accumColor    = vec3(0.0);

    // perform ray marching
    float t = tInner;
    while(t <= tOuter) {
        // get position within cloud layer
        vec3 pos = ray.o + ray.dir*t;
        float h = remap(pos.y, tInner, tOuter, 0, 1);

        // calculate density and perform alpha blending
        float d = density(pos, h, uWindDir, uTime, uGlobalCoverage, uGlobalDensity);
        //alpha += (1-alpha)*d;
        alpha += d;

        // advance ray position based on the current stepsize
        t += stepSize;

        // early ray termination
        if(alpha > 1.0) { alpha = 1.0; break; }
    }

    // calculate light color from light
    fragColor = vec4(vec3(alpha), 1.0);

    if(ray.dir.y <= 0) {
        fragColor = vec4(0.0, 0.0, 1.0, 1.0);
    }

    // debug
    //vec3 startpos = (ray.o + ray.dir*tInner) / 4000;
    //fragColor = texture(cloudMapTex, startpos.xz);
}