import { Link } from 'react-router-dom';
import { Menu, Dropdown, Button, Input } from 'antd';
const { Search } = Input;

const menu = (
  <Menu>
    <Menu.Item key="1">Option 1</Menu.Item>
    <Menu.Item key="2">Option 2</Menu.Item>
    <Menu.Item key="3">Option 3</Menu.Item>
    <Menu.Item key="4">Option 4</Menu.Item>
  </Menu>
);

const Header = () => {
  const onSearch = (value) => {
    console.log(value);
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '10px', backgroundColor: '#f0f0f0' }}>
      <Link to="/">Home</Link>
      <Search
        placeholder="Search..."
        onSearch={onSearch}
        style={{ width: 200 }}
      />
      <Dropdown overlay={menu}>
        <Button>
          Menu
        </Button>
      </Dropdown>
    </div>
  );
};

export default Header;