#version 330
layout (location = 0) in vec3 vPos;
layout (location = 1) in vec2 texCoord;
//layout (location = 1) in vec3 vColor;

//out vec3 vertexColor;
out vec2 TexCoord;

uniform mat4 model;
uniform mat4 view;
//uniform mat4 proj;

void main() {
    gl_Position = view * model * vec4(vPos, 1.0);
    //vertexColor = vec3(0.1, 0.3, 0.5);

    TexCoord = texCoord;
}