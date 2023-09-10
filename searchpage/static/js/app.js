function main() {
  let resultsContainer = document.getElementById("results");
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);

  fetch("/search/raw?q=" + urlParams.get("q")).then((r) => {
    r.json().then((r) => {
      let rString = `${r.results.length} results in ${r.time / 1000}s`;
      for (let i = 0; i < r.results.length; i++) {
        const result = r.results[i];
        rString += `<div class="result"><a href="${result.address}"><h3>Result</h3></a><p>${result.description}</p></div>`;
      }
      resultsContainer.innerHTML = rString;
    });
  });
}
