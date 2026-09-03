const data = rq.response.json();
if (data.accessToken) {
  rq.environment.set('ACCESS_TOKEN', data.accessToken);
  console.log('ACCESS_TOKEN updated successfully');
} else {
  console.error('accessToken not found in response');
}
