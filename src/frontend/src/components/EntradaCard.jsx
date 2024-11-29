import React from "react";
import { Container, Row, Col } from "react-bootstrap";
import { Link, useNavigate } from "react-router-dom";

const EntradaCard = ({ id, title, author, createdAt, onEntradaClick }) => {
  const navigate = useNavigate();

  const handleClick = () => {
    onEntradaClick(id);
    navigate(`/entrada/${id}`);
  };

  return (
    <Link onClick={handleClick} className="custom-link m-0 p-0">
      <Container className={`border border-dark-subtle bg-light text-dark mb-2 mx-0 p-0`}>
        <Row className={"custom-link-container m-0 p-0"}>
          <Col className="p-3 d-flex align-items-center">
            <span className="h4">{title}</span>
          </Col>
          <Col className="p-3">
            <Row>
              <Col className="m-auto">
                <span className="h5">Author</span>
                <p>{author}</p>
              </Col>
              <Col className="m-auto">
                <span className="h5">Created At</span>
                <p>{createdAt}</p>
              </Col>
            </Row>
          </Col>
        </Row>
      </Container>
    </Link>
  );
};

export default EntradaCard;
