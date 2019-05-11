// calculate cone sampling similar to Meteoros
// with the last sample being 3 times the 
// unit cone size
void coneSampling(inout vec3 coneSamples[5]) {
    coneSamples[0] = vec3(0.1, 0.1, 0.0);
    coneSamples[1] = vec3(0.2, 0.0, 0.2);
    coneSamples[2] = vec3(0.4,-0.4, 0.0);
    coneSamples[3] = vec3(0.8, 0.0,-0.8);
    coneSamples[4] = vec3(3.0, 0.0, 0.0);
}

// creates a rotation matrix that rotates vector a onto 
// vector b by performing a 2D rotation on the plane
// with normal a x b
// source: https://math.stackexchange.com/questions/180418/calculate-rotation-matrix-to-align-vector-a-to-vector-b-in-3d
mat3 coneRotationMatrix(in vec3 coneDir, in vec3 lightDir) {
    vec3 a = normalize(coneDir);
    vec3 b = normalize(lightDir);
    return mat3(
        dot(a, b)          , -length(cross(a, b)), 0.0,
        length(cross(a, b)),  dot(a, b)          , 0.0,
        0.0                ,  0.0                , 1.0
    );
}

// take 5 samples in direction to the sun, where the last
// sample is further away than the rest to also capture
// clouds further away
float sampleConeDensity(in vec3 pos, in vec3 coneSamples[5]) {
    // calculate orientation towards the sun
    vec3 toLight = normalize(sunPos - pos);
    mat3 coneRot = coneRotationMatrix(vec3(1, 0, 0), toLight);

    // sample density along cone
    float cloudDensity= 0.0;
    for(int i = 0; i < 5; i++) {
        vec3 coneSamplePos = pos + (coneRot * coneSamples[i] * 10);
        float coneSampleHeight = relativeHeight(coneSamplePos);
        cloudDensity += sampleCloudBase(coneSamplePos, coneSampleHeight);
    }
    return cloudDensity;
}