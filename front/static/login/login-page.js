const loginForm = document.getElementById("login-form");
const loginButton = document.getElementById("login-form-submit");
const loginErrorMsg = document.getElementById("login-error-msg");

// When the login button is clicked, the following code is executed
loginButton.addEventListener("click", (e) => {
    // Prevent the default submission of the form
    e.preventDefault();
    // Get the values input by the user in the form fields
    const username = loginForm.username.value;
    const password = loginForm.password.value;

    const params = {
        id: username,
        pass: password
    };
    const options = {
        method: 'POST',
        body: JSON.stringify( params )  
    };
    fetch( 'http://localhost:8080/login/auth', options )
        .then( response => response.text() )
        .then( response => {
            if (response == "OK") {
                location.reload();
            } else {
                loginErrorMsg.style.opacity = 1;
            }
        } );
    

})