#version 430

in VertexOut {
    vec2 uv;
}i;
out vec3 fragColor;

uniform float threshold;
uniform sampler2D tex;

float luminocity(vec3 color) {
    float lum = dot(color, vec3(0.30, 0.59, 0.11));
    // f(x) = 0,  x <= threshold
    //        x,  else
    if (lum <= threshold) {
        lum = 0.0;
    }
    return lum; // * sign(clamp(lum - threshold, 0, lum));
}

void main() {
    float lum = luminocity(texture(tex, i.uv).xyz);
    fragColor = vec3(lum, lum, lum);
}