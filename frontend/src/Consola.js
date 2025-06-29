import React, { useState } from 'react';
import './style.css';

export const Consola = ({ onGoToLogin}) => { 
    const [inputCommands, setInputCommands] = useState('');
    const [outputMessages, setOutputMessages] = useState('');

    const handleFileUpload = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                setInputCommands(e.target.result);
            };
            reader.readAsText(file);
        }
    };

    const handleExecuteCommands = async () => {
        try {
            const response = await fetch('http://localhost:8080/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ commands: inputCommands }),
            });

            if (response.ok) {
                const data = await response.json();
                setOutputMessages(data.output);
            } else {
                setOutputMessages('Error al ejecutar los comandos.');
            }
        } catch (error) {
            setOutputMessages('Error de red al intentar conectarse con el servidor.');
            console.error('Network error:', error);
        }
    };

    return (
        <div className="container">
            <h1>Ejecutor de Comandos</h1>
            <div className="flex-container">
                <div className="input-area">
                    <label htmlFor="inputCommands">Comandos a Ejecutar:</label>
                    <textarea
                        id="inputCommands"
                        placeholder="Ingresa los comandos aquí..."
                        value={inputCommands}
                        onChange={(e) => setInputCommands(e.target.value)}
                    ></textarea>
                </div>

                <div className="output-area">
                    <label htmlFor="outputMessages">Salida:</label>
                    <textarea
                        id="outputMessages"
                        readOnly
                        placeholder="Aquí se mostrarán los resultados de los comandos..."
                        value={outputMessages}
                    ></textarea>
                </div>
            </div>

            <div className="button-group">
                <input
                    type="file"
                    id="loadFile"
                    onChange={handleFileUpload}
                />
                <label htmlFor="loadFile" className="button">Cargar Archivo</label>
                <button id="executeCommands" onClick={handleExecuteCommands}>Ejecutar</button>
                <button onClick={onGoToLogin} className="button">Iniciar Seción</button> 
            </div>
        </div>
    );
};