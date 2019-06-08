#version 430
//---------------------------------------------------------------------------------------//
// textures                                                                              //
//---------------------------------------------------------------------------------------//
layout(binding = 0) uniform sampler3D cloudBaseTex;
layout(binding = 1) uniform sampler3D cloudDetailTex;
layout(binding = 2) uniform sampler2D turbulenceTex;
layout(binding = 3) uniform sampler2D cloudMapTex;

//---------------------------------------------------------------------------------------//
// includes                                                                              //
//---------------------------------------------------------------------------------------//
#include "util/camera.glsl"
#include "util/ray.glsl"
#include "util/math.glsl"
#include "cloud/density.glsl"

//---------------------------------------------------------------------------------------//
// uniforms                                                                              //
//---------------------------------------------------------------------------------------//
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

//---------------------------------------------------------------------------------------//
// input                                                                                 //
//---------------------------------------------------------------------------------------//
in Vertex {
    vec2 uv;
} i;

//---------------------------------------------------------------------------------------//
// output                                                                                //
//---------------------------------------------------------------------------------------//
out vec4 fragColor;

//---------------------------------------------------------------------------------------//
// helper functions                                                                      //
//---------------------------------------------------------------------------------------//
float intersectLayer(in Ray ray, in float h) {
    return (h - ray.o.y) / ray.dir.y;
}

//---------------------------------------------------------------------------------------//
// entry point                                                                           //
//---------------------------------------------------------------------------------------//
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
        alpha += (1-alpha)*d;

        // advance ray position based on the current stepsize
        t += stepSize;

        // early ray termination
        if(alpha > 1.0) { alpha = 1.0; break; }
    }

    // calculate light color from light
    fragColor = vec4(vec3(alpha), 1.0);

    // debug
    //vec3 startpos = (ray.o + ray.dir*tInner) / 4000;
    //fragColor = texture(cloudMapTex, startpos.xz);
}