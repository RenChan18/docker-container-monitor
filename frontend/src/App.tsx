// frontend/src/App.tsx
import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert } from 'antd';
import axios from 'axios';

interface ContainerStatus {
  id: number;
  ip_address: string;
  ping_duration: number;
  last_successful: string;
}

const App: React.FC = () => {
  const [data, setData] = useState<ContainerStatus[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    try {
      const response = await axios.get<ContainerStatus[]>('http://localhost:8080/containers');
      setData(response.data);
      setLoading(false);
    } catch (err: any) {
      setError(err.message);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000); 
    return () => clearInterval(interval);
  }, []);

  const columns = [
    {
      title: 'IP адрес',
      dataIndex: 'ip_address',
      key: 'ip_address',
    },
    {
      title: 'Ping Duration (ms)',
      dataIndex: 'ping_duration',
      key: 'ping_duration',
    },
    {
      title: 'Последний успешный пинг',
      dataIndex: 'last_successful',
      key: 'last_successful',
      render: (text: string) => new Date(text).toLocaleString(),
    },
  ];

  if (loading) return <Spin tip="Loading..." />;
  if (error) return <Alert message="Error" description={error} type="error" showIcon />;

  return (
    <div style={{ padding: '20px' }}>
      <h1>Статусы контейнеров</h1>
      <Table dataSource={data} columns={columns} rowKey="id" />
    </div>
  );
};

export default App;

