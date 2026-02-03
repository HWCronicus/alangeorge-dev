import { useEffect, useRef } from "react";
import * as THREE from "three";
import { FontLoader } from "three/examples/jsm/loaders/FontLoader.js";
import { TextGeometry } from "three/examples/jsm/geometries/TextGeometry.js";

export default function ThreeScene() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!canvasRef.current) return;

    let animationId: number;
    let hoverOffset = 0;
    let textMesh: THREE.Mesh | null = null;

    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x050505);
    scene.fog = new THREE.Fog(0x050505, 10, 50);

    const camera = new THREE.PerspectiveCamera(
      75,
      window.innerWidth / window.innerHeight,
      0.1,
      1000,
    );
    camera.position.set(0, 4, 12);
    camera.lookAt(0, 1, 0);

    const renderer = new THREE.WebGLRenderer({
      canvas: canvasRef.current,
      antialias: true,
    });
    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.useLegacyLights = true;
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;

    const ambientLight = new THREE.AmbientLight(0xffffff, 0.1);
    scene.add(ambientLight);

    const pointLight1 = new THREE.PointLight(0xff8800, 0.3, 50);
    pointLight1.position.set(10, 10, 10);
    scene.add(pointLight1);

    const pointLight2 = new THREE.PointLight(0xffaa33, 0.2, 50);
    pointLight2.position.set(-10, 8, -10);
    scene.add(pointLight2);

    const flashlight = new THREE.SpotLight(0xffffff, 2);
    flashlight.position.set(0, 10, 12);
    flashlight.angle = Math.PI / 4;
    flashlight.penumbra = 0.5;
    flashlight.decay = 1;
    flashlight.distance = 100;
    flashlight.castShadow = true;
    flashlight.shadow.mapSize.width = 2048;
    flashlight.shadow.mapSize.height = 2048;
    flashlight.shadow.camera.near = 1;
    flashlight.shadow.camera.far = 100;
    flashlight.shadow.camera.fov = 50;
    flashlight.shadow.bias = -0.0005;

    const flashlightTarget = new THREE.Object3D();
    flashlightTarget.position.set(0, 1.2, 0);
    scene.add(flashlightTarget);
    flashlight.target = flashlightTarget;
    scene.add(flashlight);

    const floorSize = 100;
    const gridSize = 50;
    const squareSize = floorSize / gridSize;
    const floorGroup = new THREE.Group();

    for (let i = 0; i < gridSize; i++) {
      for (let j = 0; j < gridSize; j++) {
        const isLight = (i + j) % 2 === 0;
        const squareGeometry = new THREE.PlaneGeometry(squareSize, squareSize);
        const squareMaterial = new THREE.MeshStandardMaterial({
          color: isLight ? 0x1a1a1a : 0x0f0f0f,
          metalness: 0.3,
          roughness: 0.8,
        });
        const square = new THREE.Mesh(squareGeometry, squareMaterial);
        square.rotation.x = -Math.PI / 2;
        square.position.set(
          (i - gridSize / 2) * squareSize + squareSize / 2,
          0,
          (j - gridSize / 2) * squareSize + squareSize / 2,
        );
        square.receiveShadow = true;
        floorGroup.add(square);
      }
    }
    scene.add(floorGroup);

    const loader = new FontLoader();
    loader.load(
      "/fonts/fira_code.json",
      (font) => {
        const textGeometry = new TextGeometry("AlanGeorge.Dev", {
          font: font,
          size: 2.0,
          height: 0.2,
          curveSegments: 12,
          bevelEnabled: true,
          bevelThickness: 0.15,
          bevelSize: 0.05,
          bevelOffset: 0,
          bevelSegments: 15,
        });

        textGeometry.computeBoundingBox();
        const centerOffset =
          -0.5 *
          (textGeometry.boundingBox!.max.x - textGeometry.boundingBox!.min.x);

        const textMaterial = new THREE.MeshStandardMaterial({
          color: 0xff8800,
          metalness: 0.5,
          roughness: 0.5,
          emissive: 0xff8800,
          emissiveIntensity: 0.1,
        });

        textMesh = new THREE.Mesh(textGeometry, textMaterial);
        textMesh.position.x = centerOffset;
        textMesh.position.y = 1.8;
        textMesh.position.z = 0;
        textMesh.castShadow = true;
        scene.add(textMesh);

        const descGeometry = new TextGeometry(
          "Full Web Stack Developer & Overall Cool Guy",
          {
            font: font,
            size: 0.5,
            height: 0.08,
            curveSegments: 8,
            bevelEnabled: true,
            bevelThickness: 0.02,
            bevelSize: 0.02,
            bevelOffset: 0,
            bevelSegments: 3,
          },
        );

        descGeometry.computeBoundingBox();
        const descCenterOffset =
          -0.5 *
          (descGeometry.boundingBox!.max.x - descGeometry.boundingBox!.min.x);

        const descMaterial = new THREE.MeshStandardMaterial({
          color: 0xffaa33,
          metalness: 0.4,
          roughness: 0.3,
          emissive: 0xff8800,
          emissiveIntensity: 0.05,
        });

        const descMesh = new THREE.Mesh(descGeometry, descMaterial);
        descMesh.position.x = descCenterOffset;
        descMesh.position.y = 0.3;
        descMesh.position.z = 0;
        descMesh.castShadow = true;
        scene.add(descMesh);
      },
      undefined,
      (error) => {
        console.error("Font loading error:", error);
      },
    );

    const handleMouseMove = (event: MouseEvent) => {
      const mouseX = (event.clientX / window.innerWidth) * 2 - 1;
      const mouseY = -(event.clientY / window.innerHeight) * 2 + 1;

      flashlight.position.x = mouseX * 15;
      flashlight.position.y = 8 + mouseY * 5;
      flashlight.position.z = 12;

      flashlightTarget.position.x = mouseX * 0.5;
      flashlightTarget.position.y = 1.2 + mouseY * 0.5;
    };
    window.addEventListener("mousemove", handleMouseMove);

    const handleResize = () => {
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(window.innerWidth, window.innerHeight);
    };
    window.addEventListener("resize", handleResize);

    const animate = () => {
      animationId = requestAnimationFrame(animate);
      hoverOffset += 0.02;

      if (textMesh) {
        textMesh.position.y = 2.0 + Math.sin(hoverOffset) * 0.15;
        textMesh.rotation.y = Math.sin(hoverOffset * 0.5) * 0.05;
      }

      renderer.render(scene, camera);
    };
    animate();

    return () => {
      cancelAnimationFrame(animationId);
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("mousemove", handleMouseMove);
      renderer.dispose();
    };
  }, []);

  return <canvas ref={canvasRef} id="canvas" />;
}
