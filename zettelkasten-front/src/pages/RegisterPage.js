import React, { useState } from 'react';
import { createUser } from "../api";

function RegisterPage() {
    // State to store each input field's value
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');

    // Handle form submission
    const handleSubmit = (e) => {
        e.preventDefault();

        // Basic validation to check if passwords match
        if (password !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        // Submit the form data
        console.log('Form submitted', { username, email, password, confirmPassword });
        // Here you would typically send the data to your server via an API

	const userData = { username, email, password, confirmPassword };

	createUser(userData)
	    .then(data => {
		console.log('User created successfully', data);
		// Handle successful user creation (e.g., redirecting the user or showing a success message)
		// Reset form or redirect user to login page, etc.
	    })
	    .catch(error => {
		console.error('Error creating user:', error);
		setError('Failed to create user. Please try again.');
	    });
    };

    return (
        <div>
            <h2>Register</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Username:</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label>Email:</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label>Password:</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <div>
                    <label>Confirm Password:</label>
                    <input
                        type="password"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                    />
                </div>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                <button type="submit">Register</button>
            </form>
        </div>
    );
}

export default RegisterPage;