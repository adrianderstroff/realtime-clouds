#version 430

const float PI        = 3.1415926535897932384626433832795;
const float PI_2      = 2*PI;
const float DEG_2_RAD = PI/180.0;

layout(binding = 0) uniform sampler3D cloudBaseTex;
layout(binding = 1) uniform sampler3D cloudDetailTex;
layout(binding = 2) uniform sampler2D turbulenceTex;
layout(binding = 3) uniform sampler2D cloudMapTex;

uniform vec3 cameraPos;
uniform mat4 M, V, P;
uniform float width, height;
uniform float fov = 45.0;
uniform vec3 windDir = vec3(1, 0, 1);

in Vertex {
    vec2 uv;
} i;

out vec4 fragColor;

// maps one value from one interval [inMin,inMax] to another interval [outMin, outMax]
float remap(in float val, in float inMin, in float inMax, in float outMin, in float outMax) {
    return (val - inMin)/(inMax - inMin) * (outMax - outMin) + outMin;
}

// swaps a and b
void swap(inout float a, inout float b) {
    float t = b;
    b = a;
    a = t;
}

// gets the ray origin and direction for the current fragment
void ray(out vec3 o, out vec3 dir) {
    // calc image plane
    float ar = width / height;
    float angle = tan(fov/2 * DEG_2_RAD);
    vec2 imagePlane = (i.uv*2 - vec2(1)) * vec2(angle) * vec2(ar, 1);

    // extract camera space
    mat3 cameraToWorld = transpose(mat3(V));

    o = cameraPos;
    //o = cameraToWorld * o;

    vec3 pLocal = dir = vec3(imagePlane, -1);
    vec3 pWorld = cameraToWorld * pLocal;
    dir = normalize(pWorld);
}

// Returns the intersection of the closer intersection point of a ray with a sphere
// @return a positive value if the sphere has been hit, while a negative value indicates no hit
// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-sphere-intersection
float intersectSphere(in vec3 oRay, in vec3 dRay, in vec3 oSphere, in float rSphere) {
    vec3 L = oRay - oSphere;

    // calculate the determinant to check for intersection
    float a = dot(dRay, dRay);
    float b = 2.0 * dot(dRay, L);
    float c = dot(L, L) - (rSphere * rSphere);
    float d = b*b - 4*a*c;

    // determine the intersection parameters t0, t1
    float t0, t1;
    if (d < 0.0) {
        return -1.0;
    } else if (d == 0) {
        t0 = -0.5 * b/a;
        t1 = t0;
    } else {
        float q = (b>0) ? 
            -0.5 * (b + sqrt(d)) :
            -0.5 * (b - sqrt(d));
        t0 = q / a;
        t1 = c / q;
    }
    // t0 should be the smaller value
    if(t0 > t1) swap(t0, t1);

    // check if t0 and or t1 are negative, if t0 is negative then use t1
    // instead and if both are negative then they are behind the ray origin
    if(t0 < 0) {
        t0 = t1;
        if(t0 < 0) return -1;
    }

    return t0;
}

vec4 sampleCloudBase(vec3 pos, float h) {
    float x = mod(pos.x, 128) / 128;
    float y = h;
    float z = mod(pos.z, 128) / 128;
    return texture(cloudBaseTex, vec3(x, y, z));
}

// https://en.wikipedia.org/wiki/UV_mapping
vec4 sampleCloudMap(in vec3 pos, in vec3 oSphere) {
    vec3 dir = normalize(oSphere - pos);
    float x = 0.5 + atan(dir.z, dir.x) / (2*PI);
    float y = 0.5 - asin(dir.y) / PI;
    if (dir.y > 0.0) return vec4(0, 0, 0, 1);
    return texture(cloudMapTex, vec2(x, y));
}

// returns the cloud density at the specified position
float density(vec3 pos, float height) {
    vec3 newPos = pos;
    
    vec4 cloudBase = sampleCloudBase(newPos, height);
    float lowFreqNoise = cloudBase.r;
    float highFreqNoise = dot(cloudBase.gba, vec3(0.625, 0.25, 0.125));
    lowFreqNoise = clamp(lowFreqNoise, highFreqNoise, 1.0);
    float baseCloud = remap(lowFreqNoise, highFreqNoise, 1.0, 0.0, 1.0);
    baseCloud = clamp(baseCloud, 0, 1);
    
    return baseCloud;
}

void main2() {
    // determine ray direction
    vec3 o, dir;
    ray(o, dir);

    // get start and end points of the horizon
    float tInner = intersectSphere(o, dir, vec3(0, 0, 0), 1400);
    float tOuter = intersectSphere(o, dir, vec3(0, 0, 0), 4000);

    // this shouldn't happen
    //if(tInner < 0 || tOuter < 0) discard;

    // step size
    float stepSize = (tOuter - tInner) / 96.0;

    // perform ray marching
    float light = 0;
    float t = tInner;
    int sampleCount = 0;
    while(t <= tOuter) {
        // get position within cloud layer
        vec3 pos = o + dir*t;
        float h = remap(t, tInner, tOuter, 0, 1);
        sampleCount++;

        // calculate density
        float d = density(pos, h);

        // do light calculations
        light += d;

        // do step
        t += stepSize;
    }
    light /= sampleCount;

    // calculate light color from light
    fragColor = vec4(light, light, light, 1);
    if(tInner < 0 || tOuter < 0) fragColor = vec4(1, 0, 0, 1);
    fragColor = sampleCloudMap(dir, vec3(0, 0, 0));
}

void main() {
    // determine ray direction
    vec3 o, dir;
    ray(o, dir);

    float t = intersectSphere(o, dir, vec3(0, 0, 0), 4000);
    fragColor = sampleCloudMap(o+t*dir, vec3(0, 0, 0));
    if(t < 0) fragColor = vec4(1, 1, 1, 1);
}