"use client";

import { useEffect, useState } from "react";
import {
    Button,
    Card,
    Descriptions,
    Empty,
    Flex,
    Form,
    Input,
    Space,
    Tag,
    Typography,
    App,
    Select,
} from "antd";
import { useParams, useRouter } from "next/navigation";
import { apiRequest } from "@/lib/api";
import { getAuthPayload, getToken, removeToken } from "@/lib/auth";
import type {
    CreateReplyRequest,
    GetTicketResponse,
    ListRepliesResponse,
    Ticket,
    TicketReply,
} from "@/types/api";
import RecruiterHint from "@/components/RecruiterHint";

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

export default function TicketDetailPage() {
    const { message } = App.useApp();
    const params = useParams<{ id: string }>();
    const router = useRouter();
    const [form] = Form.useForm<CreateReplyRequest>();

    const [role, setRole] = useState<"customer" | "agent" | "admin" | null>(null);
    const [updatingStatus, setUpdatingStatus] = useState(false);
    const [assigning, setAssigning] = useState(false);
    const [selectedStatus, setSelectedStatus] = useState<string>("");
    const [assignedTo, setAssignedTo] = useState<number | null>(null);

    const [token, setToken] = useState<string | null>(null);
    const [ticket, setTicket] = useState<Ticket | null>(null);
    const [replies, setReplies] = useState<TicketReply[]>([]);
    const [loadingTicket, setLoadingTicket] = useState(false);
    const [loadingReplies, setLoadingReplies] = useState(false);
    const [submittingReply, setSubmittingReply] = useState(false);

    const ticketId = params.id;

    useEffect(() => {
        const storedToken = getToken();
        if (!storedToken) {
            router.replace("/login");
            return;
        }

        const payload = getAuthPayload();
        setRole(payload?.role ?? null);
        setToken(storedToken);
    }, [router]);

    const fetchTicket = async (authToken: string) => {
        try {
            setLoadingTicket(true);
            const response = await apiRequest<GetTicketResponse>(`/tickets/${ticketId}`, {
                method: "GET",
                token: authToken,
            });
            setTicket(response.data);
            setSelectedStatus(response.data.status);
            setAssignedTo(response.data.assigned_to);
        } catch (error) {
            const err = error as Error;
            message.error(err.message || "Failed to fetch ticket");
        } finally {
            setLoadingTicket(false);
        }
    };

    const fetchReplies = async (authToken: string) => {
        try {
            setLoadingReplies(true);
            const response = await apiRequest<ListRepliesResponse>(
                `/tickets/${ticketId}/replies`,
                {
                    method: "GET",
                    token: authToken,
                }
            );
            setReplies(response.data ?? []);
        } catch (error) {
            const err = error as Error;
            message.error(err.message || "Failed to fetch replies");
            setReplies([]);
        } finally {
            setLoadingReplies(false);
        }
    };

    const handleUpdateStatus = async () => {
        if (!token || !selectedStatus) return;

        try {
            setUpdatingStatus(true);

            await apiRequest(`/tickets/${ticketId}/status`, {
                method: "PATCH",
                token,
                body: { status: selectedStatus },
            });

            message.success("Status updated");
            await fetchTicket(token);
        } catch (error) {
            const err = error as Error;
            message.error(err.message || "Failed to update status");
        } finally {
            setUpdatingStatus(false);
        }
    };

    const handleAssignTicket = async () => {
        if (!token || !assignedTo) return;

        try {
            setAssigning(true);

            await apiRequest(`/tickets/${ticketId}/assign`, {
                method: "PATCH",
                token,
                body: { assigned_to: assignedTo },
            });

            message.success("Ticket assigned");
            await fetchTicket(token);
        } catch (error) {
            const err = error as Error;
            message.error(err.message || "Failed to assign ticket");
        } finally {
            setAssigning(false);
        }
    };

    useEffect(() => {
        if (!token) return;

        fetchTicket(token);
        fetchReplies(token);
    }, [token, ticketId]);

    const handleLogout = () => {
        removeToken();
        router.push("/login");
    };

    const handleReplySubmit = async (values: CreateReplyRequest) => {
        if (!token) return;

        try {
            setSubmittingReply(true);

            await apiRequest(`/tickets/${ticketId}/replies`, {
                method: "POST",
                token,
                body: values,
            });

            message.success("Reply posted");
            form.resetFields();
            await fetchReplies(token);
        } catch (error) {
            const err = error as Error;
            message.error(err.message || "Failed to post reply");
        } finally {
            setSubmittingReply(false);
        }
    };

    const renderStatusTag = (status: Ticket["status"]) => {
        switch (status) {
            case "open":
                return <Tag>Open</Tag>;
            case "in_progress":
                return <Tag color="gold">In Progress</Tag>;
            case "resolved":
                return <Tag color="green">Resolved</Tag>;
            default:
                return <Tag>{status}</Tag>;
        }
    };

    const renderPriorityTag = (priority: Ticket["priority"]) => {
        switch (priority) {
            case "low":
                return <Tag>Low</Tag>;
            case "medium":
                return <Tag color="orange">Medium</Tag>;
            case "high":
                return <Tag color="red">High</Tag>;
            default:
                return <Tag>{priority}</Tag>;
        }
    };

    return (
        <div style={{ padding: 24 }}>
            <Flex justify="space-between" align="center" style={{ marginBottom: 24 }}>
                <Space>
                    <Button onClick={() => router.push("/tickets")}>Back</Button>
                    <Title level={2} style={{ margin: 0 }}>
                        Ticket Detail
                    </Title>
                </Space>

                <Button danger onClick={handleLogout}>
                    Logout
                </Button>
            </Flex>

            <div
                style={{
                    display: "grid",
                    gridTemplateColumns: "2fr 1fr",
                    gap: 24,
                    alignItems: "start",
                }}
            >
                <Space orientation="vertical" size="large" style={{ width: "100%" }}>
                    <Card loading={loadingTicket} title="Conversation">
                        {ticket && (
                            <>
                                <Title level={4}>{ticket.title}</Title>
                                <Paragraph>{ticket.description}</Paragraph>
                            </>
                        )}

                        {loadingReplies ? (
                            <Text type="secondary">Loading replies...</Text>
                        ) : replies.length === 0 ? (
                            <Empty description="No replies yet" />
                        ) : (
                            <Space orientation="vertical" size="middle" style={{ width: "100%" }}>
                                {replies.map((reply) => (
                                    <Card key={reply.id} size="small">
                                        <Flex justify="space-between" align="center">
                                            <Text strong>User #{reply.user_id}</Text>
                                            <Text type="secondary">
                                                {new Date(reply.created_at).toLocaleString()}
                                            </Text>
                                        </Flex>
                                        <Paragraph style={{ marginTop: 8, marginBottom: 0 }}>
                                            {reply.message}
                                        </Paragraph>
                                    </Card>
                                ))}
                            </Space>
                        )}
                    </Card>

                    <Card title="Add Reply">
                        <Form<CreateReplyRequest>
                            form={form}
                            layout="vertical"
                            onFinish={handleReplySubmit}
                        >
                            <Form.Item
                                label="Message"
                                name="message"
                                rules={[{ required: true, message: "Please enter a reply" }]}
                            >
                                <TextArea rows={4} placeholder="Type your reply..." />
                            </Form.Item>

                            <Form.Item style={{ marginBottom: 0 }}>
                                <Button type="primary" htmlType="submit" loading={submittingReply}>
                                    Send Reply
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </Space>

                <Card title="Ticket Info" loading={loadingTicket}>
                    {ticket && (
                        <Space orientation="vertical" size="middle" style={{ width: "100%", marginTop: 16 }}>
                            {role === "customer" && (
                                <Descriptions column={1} size="small">
                                    <Descriptions.Item label="Status">
                                        {renderStatusTag(ticket.status)}
                                    </Descriptions.Item>
                                    <Descriptions.Item label="Priority">
                                        {renderPriorityTag(ticket.priority)}
                                    </Descriptions.Item>
                                </Descriptions>
                            )}
                            {(role === "agent" || role === "admin") && (
                                <Card size="small" title="Update Status">
                                    <Space orientation="vertical" style={{ width: "100%" }}>
                                        <Select
                                            value={selectedStatus}
                                            onChange={setSelectedStatus}
                                            options={[
                                                { value: "open", label: "Open" },
                                                { value: "in_progress", label: "In Progress" },
                                                { value: "resolved", label: "Resolved" },
                                            ]}
                                        />
                                        <Button
                                            type="primary"
                                            loading={updatingStatus}
                                            onClick={handleUpdateStatus}
                                        >
                                            Update Status
                                        </Button>
                                    </Space>
                                </Card>
                            )}

                            {role === "admin" && (
                                <Card size="small" title="Assign Ticket">
                                    <Space orientation="vertical" style={{ width: "100%" }}>
                                        <Select
                                            value={assignedTo}
                                            onChange={setAssignedTo}
                                            placeholder="Select an agent"
                                            options={[
                                                { value: 2, label: "Bob (Agent)" },
                                            ]}
                                        />
                                        <Button
                                            loading={assigning}
                                            onClick={handleAssignTicket}
                                        >
                                            Assign
                                        </Button>
                                    </Space>
                                </Card>
                            )}

                            <RecruiterHint
                                text={`This page demonstrates:
- Conditional UI based on user role (customer, agent, admin)
- Loading states and error handling for all operations`}
                                tags={["..."]}
                            />
                        </Space>
                    )}
                </Card>
            </div>
        </div>
    );
}
