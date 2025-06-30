import React, { useState } from 'react';
import './style.css'; 

export const Login = ({ onLoginSuccess, goToConsola }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [partitionId, setPartitionId] = useState(''); 
    const [errorMessage, setErrorMessage] = useState('');

    const handleLogin = async () => {
        try {
            const response = await fetch('http://3.142.172.171:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password, partitionId }),
            });
    
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    onLoginSuccess(data);  // Enviar los datos al frontend cuando el login es exitoso
                } else {
                    setErrorMessage(data.message);
                }
            } else {
                setErrorMessage('Error al procesar el login.');
            }
        } catch (error) {
            setErrorMessage('Error al conectar con el servidor.');
            console.error('Network error:', error);
        }
    };
    

    return (
        <div className="login-container">
            <h1>Iniciar Sesión</h1>
            {errorMessage && <p className="error">{errorMessage}</p>}
            <div className="input-group">
                <label htmlFor="partitionId">ID Partición</label>
                <input
                    type="text"
                    id="partitionId"
                    value={partitionId}
                    onChange={(e) => setPartitionId(e.target.value)}
                    placeholder="Ingresa tu ID de Partición"
                />
            </div>
            <div className="input-group">
                <label htmlFor="username">Usuario</label>
                <input
                    type="text"
                    id="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Ingresa tu usuario"
                />
            </div>
            <div className="input-group">
                <label htmlFor="password">Contraseña</label>
                <input
                    type="password"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Ingresa tu contraseña"
                />
            </div>
            <button className='button' onClick={handleLogin}>Iniciar Sesión</button>
            <button onClick={goToConsola} className='button'>Regresar a Consola</button>
        </div>
    );
};