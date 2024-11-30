import { useEffect, useState } from "react";
import { List, Card, Alert, Typography } from "antd";
import { getAllWikis } from "../api/WikiApi.js";
import MainLayout from "../layout/MainLayout.jsx";
const { Title, Paragraph } = Typography;

function HomePage() {
  const [wikis, setWikis] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    getAllWikis()
      .then(setWikis)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <MainLayout>
      <div style={{ padding: "20px" }}>
        <Title level={1} style={{color: "white"}}>Wikis</Title>
        {error && <Alert message={error} type="error" showIcon />}
        <List
          grid={{ gutter: 16, column: 2 }}
          dataSource={wikis}
          renderItem={(wiki) => (
            <List.Item>
              <Card
                hoverable
                cover={
                  <img alt="example" src="https://via.placeholder.com/150" />
                }
              >
                <Card.Meta
                  title={<Title level={2}>{wiki.title}</Title>}
                  description={<Paragraph>{wiki.description}</Paragraph>}
                />
              </Card>
            </List.Item>
          )}
        />
      </div>
    </MainLayout>
  );
}

export default HomePage;
