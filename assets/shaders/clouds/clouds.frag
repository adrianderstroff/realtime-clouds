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
uniform vec3   cameraPos;
uniform mat4   M, V, P;
uniform float  width, height;
uniform float  fov              = 45.0;
// sun
uniform vec3   sunPos           = vec3(40000, -20000, 0);
// atmosphere
uniform float innerHeight       = 14000;
uniform float outerHeight       = 40000;
// animation
uniform float  uTime            = 0;
uniform float  uWindSpeed       = 10;
uniform vec3   windDir          = vec3(1, 0, 1);
// colors
uniform vec3   uSunColor        = vec3(1, 0, 0);
uniform vec3   uAmbientColor    = vec3(0.4, 0.4, 0.4);
uniform vec3   uAtmosphereColor = vec3(0, 0, 1);

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
#include "constants.glsl"
#include "ray.glsl"
#include "util.glsl"
#include "cloudsampling.glsl"
#include "conesampling.glsl"
#include "lighting.glsl"
#include "tonemapping.glsl"

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
    float stepSize = (tOuter - tInner) / 30.0;

    // setup cone sample
    vec3 coneSamples[5];
    coneSampling(coneSamples);

    // setup ray marching variables
    float alpha         = 0.0;
    float transmittance = 1.0;
    vec3  accumColor    = vec3(0.0);

    // perform ray marching
    float t = tInner;
    while(t <= tOuter) {
        // get position within cloud layer
        vec3 pos = ray.o + ray.dir*t;
        float h = remap(t, tInner, tOuter, 0, 1);

        // offset by time where wind speed increases over height
        pos += uWindSpeed *uTime * h;

        // calculate density and perform alpha blending
        float d = density(pos, h);
        alpha += (1-alpha)*d;

        // do light calculations and determine the color of the
        // current sample using the color of the sun
        if(d > 0) {
            float lightEnergy = radiance(pos, coneSamples, h, 0.7, 1.0);
            transmittance = mix(transmittance, lightEnergy, (1.0 - alpha));
            vec3 directLightColor = uSunColor * transmittance;
            accumColor += mix(directLightColor, uAmbientColor, clampRemap(h, 0.7, 1.0, 0.0, 1.0));
        }

        // advance ray position based on the current stepsize
        t += stepSize;

        // early ray termination
        if(alpha > 1.0) {
            alpha = 1.0;
            break;
        }
    }

    // Blend and fade out clouds into the horizon (CHANGE THIRD PARAM IN REMAP)
    const float cloudFadeOutPoint = 0.06f;
    alpha *= smoothstep(0.0, 1.0, min(1.0, remap(ray.dir.y, cloudFadeOutPoint, 0.2f, 0.0f, 1.0f)));

    // calculate light color from light
    fragColor = vec4(mix(uAtmosphereColor, accumColor, alpha), 1.0);

    // apply tone mapping
    fragColor.xyz = toneMapping(fragColor.xyz);
}