"use client";

import { useEffect, useState } from "react";
import { App, Button, Card, Form, Input, Select, Space, Typography, Flex, Grid } from "antd";
import { useRouter } from "next/navigation";
import { apiRequest } from "@/lib/api";
import { getToken, removeToken } from "@/lib/auth";
import RecruiterHint from "@/components/RecruiterHint";

const { Title } = Typography;
const { TextArea } = Input;
const { useBreakpoint } = Grid;

type CreateTicketFormValues = {
  title: string;
  description: string;
  priority: "low" | "medium" | "high";
};

export default function CreateTicketPage() {
  const { message } = App.useApp();
  const router = useRouter();
  const [form] = Form.useForm<CreateTicketFormValues>();

  const [token, setToken] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const screens = useBreakpoint();

  useEffect(() => {
    const storedToken = getToken();
    if (!storedToken) {
      router.replace("/login");
      return;
    }
    setToken(storedToken);
  }, [router]);

  const handleLogout = () => {
    removeToken();
    router.push("/login");
  };

  const handleSubmit = async (values: CreateTicketFormValues) => {
    if (!token) return;

    try {
      setSubmitting(true);

      const response = await apiRequest<{ data: { id: number } }>("/tickets", {
        method: "POST",
        token,
        body: values,
      });

      message.success("Ticket created");
      router.push(`/tickets/${response.data.id}`);
    } catch (error) {
      const err = error as Error;
      message.error(err.message || "Failed to create ticket");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ padding: screens.xs ? 12 : 24, maxWidth: 900, margin: "0 auto" }}>
      <Flex
        justify="space-between"
        align={screens.md ? "center" : "flex-start"}
        vertical={!screens.md}
        gap={12}
        style={{ marginBottom: 24 }}
      >
        <Space wrap>
          <Button onClick={() => router.push("/tickets")}>Back</Button>
          <Title level={2} style={{ margin: 0 }}>
            Create Ticket
          </Title>
        </Space>

        <Button danger onClick={handleLogout}>
          Logout
        </Button>
      </Flex>

      <Card title="New Support Request">
        <Form<CreateTicketFormValues>
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ priority: "medium" }}
        >
          <Form.Item
            label="Title"
            name="title"
            rules={[
              { required: true, message: "Please enter a title" },
              { min: 5, message: "Title must be at least 5 characters" },
            ]}
          >
            <Input placeholder="Example: Cannot log in after password reset" />
          </Form.Item>

          <Form.Item
            label="Description"
            name="description"
            rules={[
              { required: true, message: "Please enter a description" },
              { min: 10, message: "Description must be at least 10 characters" },
            ]}
          >
            <TextArea rows={6} placeholder="Describe the issue in detail..." />
          </Form.Item>

          <Form.Item
            label="Priority"
            name="priority"
            rules={[{ required: true, message: "Please select a priority" }]}
          >
            <Select
              options={[
                { value: "low", label: "Low" },
                { value: "medium", label: "Medium" },
                { value: "high", label: "High" },
              ]}
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Space wrap>
              <Button onClick={() => router.push("/tickets")}>Cancel</Button>
              <Button type="primary" htmlType="submit" loading={submitting}>
                Create Ticket
              </Button>
            </Space>
          </Form.Item>
        </Form>

        <RecruiterHint
          text={`...`}
          tags={["..."]}
        />
      </Card>
    </div>
  );
}
