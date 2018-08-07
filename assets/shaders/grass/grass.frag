#version 430

in VertexOut {
    vec3  position;
    vec2  uv;
    vec3  normal;
    float texID;
} i;

// phong shader uniforms
uniform vec3  cameraPos;
uniform vec3  lightDir;
uniform vec3  lightColor;
uniform float ambientIntensity;
uniform float diffuseIntensity;
uniform float d1;
uniform float d2;

// textures
layout(binding = 0) uniform sampler2D grassAlpha;
layout(binding = 1) uniform sampler2D grassDiffuse0;
layout(binding = 2) uniform sampler2D grassDiffuse1;
layout(binding = 3) uniform sampler2D grassDiffuse2;
layout(binding = 4) uniform sampler2D grassDiffuse3;

layout(location = 0) out vec3 fragColor;

void main() {
    // discard transparent pixels
    if(texture(grassAlpha, i.uv).x < 0.5 && i.texID != 0) {
        discard;
    }

    // select right grass texture
    vec3 grassColor = texture(grassDiffuse0, i.uv).xyz;
    float specularFactor = 0.1;
    float mixFactor = 1.0;
    if (i.texID == 3) {
        grassColor = texture(grassDiffuse3, i.uv).xyz;
        specularFactor = clamp(0.3 - i.uv.y, 0.0, 1.0)*3;
        mixFactor = 1-i.uv.y;
    } else if (i.texID == 2) {
        grassColor = texture(grassDiffuse2, i.uv).xyz;
        specularFactor = clamp(0.3 - i.uv.y, 0.0, 1.0)*3;
        mixFactor = 1-i.uv.y;
    } else if (i.texID == 1) {
        grassColor = texture(grassDiffuse1, i.uv).xyz;
        specularFactor = clamp(0.3 - i.uv.y, 0.0, 1.0)*3;
        mixFactor = 1-i.uv.y;
    }

    // combine everything
    fragColor = grassColor + (lightColor * specularFactor);
    fragColor = mix(vec3(0.0, 0.0, 0.0), fragColor, mixFactor*mixFactor);
}