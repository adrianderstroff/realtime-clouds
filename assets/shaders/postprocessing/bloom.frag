#version 430

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

layout(binding = 0) uniform sampler2D colorTexture;
layout(binding = 1) uniform sampler2D luminocityTexture;

void main() {
    vec3 color      = texture(colorTexture, i.uv).xyz;
    vec3 luminocity = texture(luminocityTexture, i.uv).xyz;
    fragColor = color + luminocity;
}