import { createSignal, onCleanup, onMount } from "solid-js";
import { createNoise3D } from "simplex-noise";
import { cn } from "@/lib/utils";

export const WavyBackground = (props) => {
  const {
    children,
    className,
    containerClassName,
    colors,
    waveWidth,
    backgroundFill,
    blur = 10,
    speed = "fast",
    waveOpacity = 0.5,
    ...rest
  } = props;

  const noise = createNoise3D();
  let w, h, nt, i, x, ctx, canvas;
  let animationId;

  let canvasRef = {};

  const getSpeed = () => {
    switch (speed) {
      case "slow":
        return 0.001;
      case "fast":
        return 0.002;
      default:
        return 0.001;
    }
  };

  const init = () => {
    canvas = canvasRef;
    ctx = canvas.getContext("2d");
    w = ctx.canvas.width = window.innerWidth;
    h = ctx.canvas.height = window.innerHeight;
    ctx.filter = `blur(${blur}px)`;
    nt = 0;

    window.onresize = () => {
      w = ctx.canvas.width = window.innerWidth;
      h = ctx.canvas.height = window.innerHeight;
      ctx.filter = `blur(${blur}px)`;
    };

    render();
  };

  const waveColors = colors || [
    "#38bdf8",
    "#818cf8",
    "#c084fc",
    "#e879f9",
    "#22d3ee",
  ];

  const drawWave = (n) => {
    nt += getSpeed();
    for (i = 0; i < n; i++) {
      ctx.beginPath();
      ctx.lineWidth = waveWidth || 50;
      ctx.strokeStyle = waveColors[i % waveColors.length];
      for (x = 0; x < w; x += 5) {
        const y = noise(x / 800, 0.3 * i, nt) * 100;
        ctx.lineTo(x, y + h * 0.5); // Adjust for height, currently at 50% of the container
      }
      ctx.stroke();
      ctx.closePath();
    }
  };

  const render = () => {
    ctx.fillStyle = backgroundFill || "black";
    ctx.globalAlpha = waveOpacity || 0.5;
    ctx.fillRect(0, 0, w, h);
    drawWave(5);
    animationId = requestAnimationFrame(render);
  };

  onMount(() => {
    init();

    onCleanup(() => {
      cancelAnimationFrame(animationId);
    });
  });

  const [isSafari, setIsSafari] = createSignal(false);

  onMount(() => {
    setIsSafari(
      typeof window !== "undefined" &&
        navigator.userAgent.includes("Safari") &&
        !navigator.userAgent.includes("Chrome")
    );
  });

  return (
    <div
      class={cn(
        "h-screen flex flex-col items-center justify-center",
        containerClassName
      )}
    >
      <canvas
        class="absolute inset-0 z-0"
        ref={(el) => (canvasRef = el)}
        id="canvas"
        style={isSafari() ? { filter: `blur(${blur}px)` } : {}}
      ></canvas>
      <div class={cn("relative z-10", className)} {...rest}>
        {children}
      </div>
    </div>
  );
};
