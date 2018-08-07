#version 430

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

layout(binding = 0) uniform sampler2D colorTexture;
layout(binding = 1) uniform sampler2D depthTexture;
layout(binding = 2) uniform sampler2D blurTexture;

uniform float zNear;
uniform float zFar;
uniform float focal;
uniform float range;

float linearizeDepth(float depth)
{
    return (2.0 * zNear) / (zFar + zNear - depth * (zFar - zNear));
}

float worldDepth(float depth) {
    return depth*(zFar - zNear) + zNear;
}

float blurIntensity(float depth) {
    return clamp(abs(depth - focal) / range, 0, 1);
}

void main() {
    vec3  color = texture(colorTexture, i.uv).xyz;
    float depth = texture(depthTexture, i.uv).r;
    vec3  blur  = texture(blurTexture,  i.uv).xyz;
    float d = blurIntensity(worldDepth(linearizeDepth(depth)));
    fragColor = color + d*(blur - color);
}