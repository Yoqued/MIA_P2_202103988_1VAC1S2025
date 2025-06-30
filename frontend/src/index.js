import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';
import { Consola } from './Consola';
import { Login } from './Login';
import { Menu } from './Menu';
import { Explorador } from './Explorador';

function Index() {
    const [isAuthenticated, setIsAuthenticated] = useState(true);
    const [activeComponent, setActiveComponent] = useState('consola');
    const [errorMessage, setErrorMessage] = useState('');

    const handleLoginSuccess = (data) => {
        setIsAuthenticated(true);
        setErrorMessage('');
        setActiveComponent('menu');
    };

    const handleLogout = async () => {
        try {
            const response = await fetch('http://3.142.172.171:8080/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ logout: true }),
            });
    
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    setIsAuthenticated(false);
                    setErrorMessage('');
                } else {
                    setErrorMessage('Error al procesar el logout.');
                }
            } else {
                setErrorMessage('Error al procesar el logout.');
            }
        } catch (error) {
            setErrorMessage('Error al conectar con el servidor.');
            console.error('Network error:', error);
        }
        
        setIsAuthenticated(false);
        setActiveComponent('login');
    };

    const goToConsola = () => {
        setIsAuthenticated(true);
        setActiveComponent('consola');
    };

    const goToFileExplorer = () => {
        setActiveComponent('fileExplorer');
    };

    const backToMenu = () => {
        setActiveComponent('menu');
    };

    const goToLogin = () => {
        setIsAuthenticated(false);
        setActiveComponent('login');
    };

    return (
        <div className="App">
            {errorMessage && <p className="error">{errorMessage}</p>}
            
            {!isAuthenticated ? (
                <Login onLoginSuccess={handleLoginSuccess} goToConsola={goToConsola} />
            ) : (
                activeComponent === 'menu' ? (
                    <Menu 
                        onLogout={handleLogout}
                        goToConsola={goToConsola}
                        goToFileExplorer={goToFileExplorer}
                    />
                ) : activeComponent === 'consola' ? (
                    <Consola onGoToLogin={goToLogin} />
                ) : activeComponent === 'fileExplorer' ? (
                    <Explorador onBackToMenu={backToMenu} />
                ) : (
                    <div>Explorador de Archivos (Próximamente)</div>
                )
            )}
        </div>
    );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<Index />);
