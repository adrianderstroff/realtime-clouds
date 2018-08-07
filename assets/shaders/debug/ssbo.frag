#version 430

const float PI = 3.14159265358;

layout(std430, binding = 0) buffer Accelerationfield { vec4 acceleration[]; };

uniform int width;
uniform int height;
uniform int ssboWidth;
uniform int ssboHeight;

out vec3 fragColor;

// transforms a HSV into a RGB value
// the input is a vec3 with all values ranging from 0 to 1
vec3 hsv2rgb(vec3 c) {
  vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
  vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
  return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

vec3 mapAccToCol(int idx) {
    vec2 acc = acceleration[idx].xy;
    vec2 dir = normalize(acc);
    float mag = length(acc);
    float h = (atan(dir.y, dir.x)) / (2*PI);
    vec3 hsv = vec3(h, 1, mag);
    return hsv2rgb(hsv);
}

void main(){
    int x = int((gl_FragCoord.x / width) * ssboWidth);
    int y = int((gl_FragCoord.y / height) * ssboHeight);
    int idx = y*ssboWidth + x;
    fragColor = mapAccToCol(idx);
}