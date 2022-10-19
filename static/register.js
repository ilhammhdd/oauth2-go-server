const xhttp = new XMLHttpRequest();

document.getElementById('submit-btn').onclick = function () {
  let csrfToken = document.getElementById('csrf-token').getAttribute('content');
  let emailAddress = document.getElementById('email_address').value;
  let username = document.getElementById('username').value;
  let password = document.getElementById('password').value;
  let confirmPassword = document.getElementById('confirm_password').value;
  if (!emailAddress && !password && !confirmPassword) {
    alert('all fields required!');
  } else if (password != confirmPassword) {
    alert('confirm password not matched!');
  } else {
    let requestPayload = {
      email_address: emailAddress,
      username: username,
      password: password,
      confirm_password: confirmPassword
    };
    console.log(requestPayload);
    // TODO: put the client_id in the url query param below
    xhttp.open('POST', '/register');
    xhttp.setRequestHeader('Csrf-Token', csrfToken);
    xhttp.send(JSON.stringify(requestPayload));
  }
}