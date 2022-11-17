document.getElementById('submit-btn').onclick = function () {
  let csrfToken = document.getElementById('csrf-token').getAttribute('content');
  let clientID = document.getElementById('client-id').getAttribute('content');
  let emailOrUsername = document.getElementById('email-address-or-username').value;
  let password = document.getElementById('password').value;
  if (!emailAddress && !password && !confirmPassword) {
    alert('all fields required!');
  } else {
    fetch("/login?client_id=" + clientID, {
      method: 'POST',
      headers: { 'Csrf-Token': csrfToken },
      body: {
        email_address_or_username: emailOrUsername,
        password: password
      }
    })
  }
}