<script>
  export let structure;
  let canvas;

  // WebGL rendering will be implemented here
  // For now, showing placeholder
</script>

<div class="viewer">
  <canvas bind:this={canvas} width="800" height="600"></canvas>
  <div class="info">
    <p>3D Protein Visualization</p>
    <p class="detail">WebGL Renderer - 60fps target</p>
  </div>
</div>

<style>
  .viewer {
    position: relative;
    background: #000;
    border-radius: 8px;
    overflow: hidden;
  }

  canvas {
    width: 100%;
    height: auto;
    display: block;
  }

  .info {
    position: absolute;
    bottom: 1rem;
    left: 1rem;
    color: #fff;
    text-shadow: 0 2px 4px rgba(0,0,0,0.8);
  }

  .detail {
    font-size: 0.85rem;
    opacity: 0.7;
  }
</style>
