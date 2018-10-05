#version 430

layout(binding = 0) uniform sampler2D rayStartTex;
layout(binding = 1) uniform sampler2D rayEndTex;
layout(binding = 2) uniform sampler3D noiseTex;
uniform int iterations = 10;

in Vertex {
    vec2 uv;
} i;

out vec4 color;

bool inside(vec3 pos, vec3 start, vec3 end) {
    return length(end - start) - length(pos - start) >= 0.0; 
}

void main() {
    vec3 start = texture(rayStartTex, i.uv).xyz;
    vec3 end   = texture(rayEndTex,   i.uv).xyz;
    vec3 dir   = normalize(end - start);

    // specify the step vector
    float stepSize = 1.0 / float(iterations);
    vec3 step = dir * stepSize;

    // the color that is accumulated during raymarching
    vec4 dst = vec4(0, 0, 0, 0);

    // do the actual raymarching
    vec3 pos = start;
    while(inside(pos, start, end)) {
        // calculate source
        vec4 src = vec4(texture(noiseTex, pos).r);
        src.a *= 0.5;
        src.rgb *= src.a;
        dst = (1.0 - dst.a) * src + dst;

        // early ray termination
        if(dst.a >= 0.99) break;

        // advance
        pos += step;
    }

    color = dst;
}