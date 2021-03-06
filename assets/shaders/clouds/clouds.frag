#version 430
//---------------------------------------------------------------------------------------//
// textures                                                                              //
//---------------------------------------------------------------------------------------//
layout(binding = 0) uniform sampler3D cloudBaseTex;
layout(binding = 1) uniform sampler3D cloudDetailTex;
layout(binding = 2) uniform sampler2D turbulenceTex;
layout(binding = 3) uniform sampler2D cloudMapTex;

//---------------------------------------------------------------------------------------//
// uniforms                                                                              //
//---------------------------------------------------------------------------------------//
// ray direction
uniform vec3  uCameraPos;
uniform mat4  M, V, P;
uniform float width, height;
uniform float fov              = 45.0;
// sun
uniform vec3  uSunPos           = vec3(40000, -1000, 0);
// atmosphere
uniform float innerHeight      = 14000;
uniform float outerHeight      = 40000;
uniform float extinctionCoeff  = 1.0/26000.0;
// animation
uniform float uTime            = 0;
uniform float uWindSpeed       = 10;
uniform vec3  windDir          = vec3(1, 0, 1);
// colors
uniform vec3  uSunColor        = vec3(1, 1, 0);
uniform vec3  uAmbientColor    = vec3(1, 0, 0);
uniform vec3  uAtmosphereColor = vec3(0.6, 0.7, 0.95);

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
// includes                                                                              //
//---------------------------------------------------------------------------------------//
#include "util/constants.glsl"
#include "util/math.glsl"
#include "util/ray.glsl"
#include "density/cloudsampling.glsl"
#include "density/conesampling.glsl"
#include "lighting/lighting.glsl"
#include "lighting/tonemapping.glsl"

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
    Ray ray = calcRay();

    // get start and end points of the horizon
    float tInner = intersectLayer(ray, innerHeight);
    float tOuter = intersectLayer(ray, outerHeight);

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
        float d = density(pos, h);
        alpha += (1-alpha)*d;

        // do light calculations and determine the color of the
        if(d > 0) {
            transmittance *= exp(-extinctionCoeff * stepSize);
            accumColor += radianceSimple(pos, h, 0.7, 1.0) * transmittance;
        }

        // advance ray position based on the current stepsize
        t += stepSize;

        // early ray termination
        if(alpha > 1.0) { alpha = 1.0; break; }
    }

    // Blend and fade out clouds into the horizon (CHANGE THIRD PARAM IN REMAP)
    const float cloudFadeOutPoint = 0.06f;
    alpha *= smoothstep(0.0, 1.0, min(1.0, remap(ray.dir.y, cloudFadeOutPoint, 0.2f, 0.0f, 1.0f)));

    // calculate light color from light
    fragColor = vec4(mix(uAtmosphereColor, accumColor, alpha), 1.0);

    // tone mapping
    //fragColor.xyz = toneMapping(fragColor.xyz);
}