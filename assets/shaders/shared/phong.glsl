struct PhongMaterial {
    vec3  ambientColor;
    float ambientIntensity;
    vec3  diffuseColor;
    float diffuseIntensity;
    vec3  specularColor;
    float specularIntensity;
    float specularPower;
};

vec3 phong(in vec3 pos, in vec3 normal, in vec3 cameraPos, in vec3 lightPos, in PhongMaterial mat) {
    vec3 toLight = normalize(lightPos - pos);
    vec3 toCamera = normalize(cameraPos - pos);

    // ambient color
    vec3 ambientColor = mat.ambientColor * mat.ambientIntensity;
    
    // diffuse color
    float diffuseFactor = dot(normal, toLight);
    vec3 diffuseColor = mat.diffuseColor * mat.diffuseIntensity * diffuseFactor;

    // specular color
    vec3 lightReflect = normalize(reflect(toLight, normal));
    float specularFactor = dot(toCamera, lightReflect);
    specularFactor = pow(specularFactor, mat.specularPower);
    vec3 specularColor = mat.specularColor * mat.specularIntensity * specularFactor;

    // put everything together
    return ambientColor + diffuseColor + specularColor;
}