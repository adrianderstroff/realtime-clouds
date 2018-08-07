#version 430

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

layout(binding = 0) uniform sampler2D colorTexture;
layout(binding = 1) uniform sampler2D depthTexture;

uniform float zNear;
uniform float zFar;
uniform vec3 cameraPos;
uniform mat4 InvViewProj;
uniform vec3 lightDir;
uniform float lightIntensity;
uniform float fogDensity;

const vec3 blue = vec3(0.5, 0.6, 0.7);
const vec3 yellow = vec3(1.0, 0.9, 0.7);

vec3 reconstructPos(float depth){
    vec4 sPos = vec4(i.uv * 2.0 - 1.0, depth, 1.0);
    sPos = InvViewProj * sPos;

    return (sPos.xyz / sPos.w);
}

vec3 applyFog(vec3 rgb, float dist, vec3 rayDir, vec3 sunDir) {
    float fogAmount = 1.0 - exp(-dist * fogDensity);
    float sunAmount = max(dot(rayDir, sunDir), 0.0);
    vec3 fogColor = mix(blue, yellow, pow(sunAmount, lightIntensity));
    return mix(rgb, fogColor, fogAmount);
}

void main() {
    vec3  color = texture(colorTexture, i.uv).xyz;
    float depth = texture(depthTexture, i.uv).r;
    
    // calc fog
    vec3 worldPos = reconstructPos(depth);
    vec3 rayDir = worldPos - cameraPos;
    float dist = length(rayDir) / (zFar-zNear);
    rayDir = normalize(rayDir);
    vec3 sunDir = -normalize(lightDir);
    fragColor   = applyFog(color, dist, rayDir, sunDir);
}