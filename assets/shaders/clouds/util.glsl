// maps one value from one interval [inMin,inMax] to another interval [outMin, outMax]
float remap(in float val, in float inMin, in float inMax, in float outMin, in float outMax) {
    return (val - inMin)/(inMax - inMin) * (outMax - outMin) + outMin;
}

// clamps the input value to (inMin, inMax) and performs a remap
float clampRemap(in float val, in float inMin, in float inMax, in float outMin, in float outMax) {
    float clVal = clamp(val, inMin, inMax);
    return (clVal - inMin)/(inMax - inMin) * (outMax - outMin) + outMin;
}

// swaps a and b
void swap(inout float a, inout float b) {
    float t = b;
    b = a;
    a = t;
}

// calculates three weights (w1, w2, w3) for a parameter t in [0, 1]
vec3 lerp3(float t) {
    float x = clamp(1 - t*2, 0, 1);
    float z = clamp((t-0.5)*2, 0, 1);
    float y = 1 - x - z;
    return vec3(x, y, z);
}