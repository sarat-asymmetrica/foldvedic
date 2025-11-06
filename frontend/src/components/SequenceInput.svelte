<script>
  import { createEventDispatcher } from 'svelte';

  export let sequence;
  export let loading = false;

  const dispatch = createEventDispatcher();

  function handleFold() {
    if (sequence.length > 0 && !loading) {
      dispatch('fold');
    }
  }
</script>

<div class="input-section">
  <label for="sequence">Protein Sequence (FASTA)</label>
  <textarea
    id="sequence"
    bind:value={sequence}
    placeholder="Enter amino acid sequence... (e.g., MKWVTFISLLLLFSSAYS)"
    rows="4"
    disabled={loading}
  ></textarea>

  <button on:click={handleFold} disabled={loading}>
    {#if loading}
      ⏳ Folding...
    {:else}
      🧬 Predict Structure
    {/if}
  </button>

  <p class="hint">
    Or upload PDB file for comparison
  </p>
</div>

<style>
  .input-section {
    background: #1a1a1a;
    border-radius: 12px;
    padding: 2rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #ddd;
    font-weight: 500;
  }

  textarea {
    width: 100%;
    padding: 1rem;
    background: #0a0a0a;
    border: 2px solid #333;
    border-radius: 8px;
    color: #fff;
    font-family: 'Courier New', monospace;
    font-size: 1rem;
    resize: vertical;
    transition: border-color 0.2s;
  }

  textarea:focus {
    outline: none;
    border-color: #667eea;
  }

  button {
    width: 100%;
    margin-top: 1rem;
    padding: 1rem 2rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    border-radius: 8px;
    color: white;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s, opacity 0.2s;
  }

  button:hover:not(:disabled) {
    transform: translateY(-2px);
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .hint {
    margin-top: 1rem;
    color: #666;
    font-size: 0.9rem;
  }
</style>
