"use client";
import { Button, Header } from "ui";
import { useState, useEffect } from "react";

interface Data {
  status: string;
  message: string;
}

export default function Page() {
  const [data, setData] = useState<Data>({ status: "", message: "" });
  // Test the private route to check if cors is properly configured
  useEffect(() => {
    fetch("http://localhost:3333/v1/private")
      .then((response) => response.json())
      .then((data) => setData(data));
  }, []);

  return (
    <>
      <Header text="Web" />
      <Button />
      <p>{data.status}</p>
      <p>{data.message}</p>
    </>
  );
}