// Fragment shader for protein atoms
// Phong shading for realistic sphere appearance

precision mediump float;

varying vec4 vColor;
varying vec3 vNormal;

uniform vec3 uLightDirection;

void main() {
    // Phong shading
    vec3 normal = normalize(vNormal);
    vec3 lightDir = normalize(uLightDirection);

    // Ambient + diffuse lighting
    float ambient = 0.3;
    float diffuse = max(dot(normal, lightDir), 0.0);
    float lighting = ambient + 0.7 * diffuse;

    vec3 finalColor = vColor.rgb * lighting;
    gl_FragColor = vec4(finalColor, vColor.a);
}
