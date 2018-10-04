#version 410 core

uniform sampler2D rayStartTex;
uniform sampler2D rayEndTex;
uniform sampler3D noiseTex;
uniform int iterations = 10;

in Vertex {
    vec2 uv;
} i;

out vec3 color;

void main() {
    vec3 start = texture(rayStartTex, i.uv).xyz;
    vec3 end   = texture(rayEndTex,   i.uv).xyz;
    vec3 dir   = normalize(end - start);

    /* float stepSize = 1.0 / float(iterations);
    vec3 step = dir * stepSize;

    vec3 pos = start;
    vec4 dst = vec4(0, 0, 0, 0);

    for(int i = 0; i < iterations; i++) {
        float value = texture(noiseTex, pos).r;

        // calculate source
        vec4 src = vec4(value);
        //src.a *= 0.5;
        src.rgb *= src.a;
        dst = (1.0 - dst.a) * src + dst;

        // early ray termination
        if(dst.a >= 0.95) break;

        // advance
        pos += step;

        //if (pos.x > 1.0 || pos.y > 1.0 || pos.z > 1.0) break;
    }

    color = dst.rgb; */
    color = texture(rayStartTex, i.uv).xyz;
}