export async function fetchUrlStatus(url) {
  try {
    const response = await fetch("http://localhost:8080/proxy?url=" + encodeURIComponent(url), { method: 'GET' });
    if (!response.ok) {
      return 'error';
    }
    const jsonBody = await response.json();
    return jsonBody.status === 'ok' ? 'ok' : 'error';
  } catch (e) {
    return 'error';
  }
}