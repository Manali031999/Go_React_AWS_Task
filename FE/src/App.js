import React, { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import GraphData from './GraphData';


function App() {
  const [instances, setInstances] = useState([]);
  const [selectedInstance, setSelectedInstance] = useState(null);
  const [formattedData, setFormattedData] = useState(null);


  useEffect(() => {
    fetchInstances();
  }, []);

  const fetchInstances = async () => {
    try {
      const response = await fetch('http://localhost:8080/instances');
      if (!response.ok) {
        throw new Error('Failed to fetch instances');
      }
      const data = await response.json();
      setInstances(data);
    } catch (error) {
      console.error('Error fetching instances:', error);
    }
  };

  const handleInstanceClick = async (instanceID) => {
    console.log('instancesID', instanceID)
    try {
      const response = await fetch(`http://localhost:8080/instances/${instanceID}`);
      if (!response.ok) {
        throw new Error(`Failed to fetch data for instance ${instanceID}`);
      }
      const data = await response.json();
      const timestamp = Object.keys(data.GraphData);
      const value = Object.values(data.GraphData);

      setSelectedInstance(data);

      var graphRes = null
      if (value.length > 0) {
        graphRes = value.map((value, index) => ({
          timestamp: timestamp[index],
          value: value,
        }))
      }
      else {
        graphRes = null
      }
      setFormattedData(graphRes)
    } catch (error) {
      console.error(`Error fetching data for instance ${instanceID}:`, error);
    }
  };


  const refreshPage = () => {
    window.location.reload(false);
  }


  return (
    <div class="App container-fluid">
      {!selectedInstance && (
        <div>
          <h1>Instances List</h1>
          <table class="table table-bordered">
            <thead class="thead-dark">
              <tr>
                <th scope="col">InstanceID</th>
                <th scope="col">Type</th>
                <th scope="col">Region</th>
              </tr>
            </thead>
            <tbody>
              {instances.map(instance => (
                <tr key={instance.InstanceID}>
                  <td onClick={() => handleInstanceClick(instance.InstanceID)}>
                    <a href="#">
                      {instance.InstanceID}
                    </a>
                  </td>
                  <td>{instance.InstanceType}</td>
                  <td>{instance.Region}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {selectedInstance && (<div>
        <h1>Selected Instance Details</h1>

        <p>ID: {selectedInstance.InstanceID}</p>
        <GraphData data={formattedData} />
        <button type="button" class="btn btn-primary" onClick={refreshPage}>Home Page!</button>
      </div>)}

    </div>
  );
}

export default App;
