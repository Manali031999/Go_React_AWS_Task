import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import moment from 'moment';

const GraphData = ({ data }) => {
    if (!data) {
        return <div>No data available</div>;
    }

    const formattedData = data.map(item => ({
        value: item.value,
        timestamp: moment(item.timestamp).format('HH:mm'),
    }));

    return (
        <div>
            <h2>Line Chart</h2>
            <LineChart width={800} height={400} data={formattedData}>
                <XAxis dataKey="timestamp" />
                <YAxis />
                <CartesianGrid stroke="#ccc" />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="value" stroke="#8884d8" />
            </LineChart>
        </div>
    );
};

export default GraphData;
