document.getElementById('submit-btn').onclick = function () {
  let csrfToken = document.getElementById('csrf-token').getAttribute('content');
  let clientID = document.getElementById('client-id').getAttribute('content');
  let emailAddress = document.getElementById('email_address').value;
  let username = document.getElementById('username').value;
  let password = document.getElementById('password').value;
  let confirmPassword = document.getElementById('confirm_password').value;
  if (!emailAddress && !password && !confirmPassword) {
    alert('all fields required!');
  } else if (password != confirmPassword) {
    alert('confirm password not matched!');
  } else {
    fetch("/register?client_id=" + clientID, {
      method: 'POST',
      headers: { 'Csrf-Token': csrfToken },
      body: JSON.stringify({
        email_address: emailAddress,
        username: username,
        password: password,
        confirm_password: confirmPassword
      })
    }).then((response) => {
      response.json().then(parsed => {
        console.log(parsed);
        if (response.status >= 200 && response.status <= 399) {
          alert('register user success');
        } else {
          alert(JSON.stringify(parsed));
        }
      })
    })
  }
}