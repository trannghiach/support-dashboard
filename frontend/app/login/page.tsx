"use client";

import { useState } from "react";
import { Button, Card, Form, Input, Typography, App } from "antd";
import { useRouter } from "next/navigation";
import { apiRequest } from "@/lib/api";
import { saveToken } from "@/lib/auth";
import type { LoginRequest, LoginResponse } from "@/types/api";

const { Title, Text } = Typography;

export default function LoginPage() {
  const router = useRouter();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm<LoginRequest>();

  const onFinish = async (values: LoginRequest) => {
    try {
      setLoading(true);

      const response = await apiRequest<LoginResponse>("/auth/login", {
        method: "POST",
        body: values,
      });

      saveToken(response.data.token);
      message.success("Login successful");
      router.push("/tickets");
    } catch (error) {
      const err = error as Error;
      message.error(err.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "grid",
        placeItems: "center",
        padding: 24,
        background: "#f5f5f5",
      }}
    >
      <Card style={{ width: 420 }}>
        <div style={{ marginBottom: 24 }}>
          <Title level={2} style={{ marginBottom: 8 }}>
            Support Dashboard
          </Title>
          <Text type="secondary">Sign in to continue</Text>
        </div>

        <Form<LoginRequest>
          form={form}
          layout="vertical"
          onFinish={onFinish}
          autoComplete="off"
        >
          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: "Please enter your email" },
              { type: "email", message: "Please enter a valid email" },
            ]}
          >
            <Input placeholder="alice@test.com" />
          </Form.Item>

          <Form.Item
            label="Password"
            name="password"
            rules={[{ required: true, message: "Please enter your password" }]}
          >
            <Input.Password placeholder="Enter password" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" block loading={loading}>
              Sign in
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}