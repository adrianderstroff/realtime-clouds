#version 430

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

uniform float sampleSize;
uniform sampler2D tex;

void main() {
    fragColor = texture(tex, i.uv*sampleSize).xyz;
}