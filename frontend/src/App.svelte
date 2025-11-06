<script>
  // FoldVedic.ai - Main Application
  import ProteinViewer from './components/ProteinViewer.svelte';
  import SequenceInput from './components/SequenceInput.svelte';

  let sequence = 'MKWVTFISLLLLFSSAYS'; // Example: Albumin signal peptide
  let prediction = null;
  let loading = false;

  async function foldProtein() {
    loading = true;

    try {
      // Call backend API
      const response = await fetch('/api/fold', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sequence })
      });

      prediction = await response.json();
    } catch (error) {
      console.error('Folding failed:', error);
    } finally {
      loading = false;
    }
  }
</script>

<main>
  <header>
    <h1>🧬 FoldVedic.ai</h1>
    <p>Protein Folding via Vedic Mathematics & Quaternion Geometry</p>
  </header>

  <div class="container">
    <SequenceInput bind:sequence on:fold={foldProtein} {loading} />

    {#if prediction}
      <div class="results">
        <h2>Prediction Results</h2>
        <ProteinViewer structure={prediction.structure} />

        <div class="metrics">
          <div class="metric">
            <span class="label">Energy:</span>
            <span class="value">{prediction.energy.toFixed(2)} kcal/mol</span>
          </div>
          <div class="metric">
            <span class="label">Vedic Score:</span>
            <span class="value">{(prediction.vedicScore * 100).toFixed(1)}%</span>
          </div>
          {#if prediction.rmsd}
            <div class="metric">
              <span class="label">RMSD:</span>
              <span class="value">{prediction.rmsd.toFixed(2)} Å</span>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #0a0a0a;
    color: #ffffff;
  }

  main {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }

  header {
    text-align: center;
    margin-bottom: 3rem;
  }

  h1 {
    font-size: 3rem;
    font-weight: 700;
    margin: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  header p {
    color: #888;
    margin-top: 0.5rem;
  }

  .container {
    display: grid;
    gap: 2rem;
  }

  .results {
    background: #1a1a1a;
    border-radius: 12px;
    padding: 2rem;
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .metric {
    background: #0a0a0a;
    padding: 1rem;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .label {
    color: #888;
    font-size: 0.9rem;
  }

  .value {
    font-size: 1.2rem;
    font-weight: 600;
    color: #667eea;
  }
</style>
