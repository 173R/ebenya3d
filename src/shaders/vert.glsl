#version 330
layout (location = 0) in vec3 vPos;
layout (location = 1) in vec3 vColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(vPos, 1.0);
    vertexColor = vColor;
}