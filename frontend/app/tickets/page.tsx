"use client";

import { useEffect, useMemo, useState } from "react";
import {
    Button,
    Card,
    Flex,
    Select,
    Space,
    Table,
    Tag,
    Typography,
    App,
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { useRouter } from "next/navigation";
import { apiRequest } from "@/lib/api";
import { getToken, removeToken } from "@/lib/auth";
import type { ListTicketsResponse, Ticket } from "@/types/api";

const { Title, Text } = Typography;

type StatusFilter = "" | "open" | "in_progress" | "resolved";
type PriorityFilter = "" | "low" | "medium" | "high";

export default function TicketsPage() {
    const { message } = App.useApp();
    const router = useRouter();

    const [tickets, setTickets] = useState<Ticket[]>([]);
    const [loading, setLoading] = useState(false);

    const [status, setStatus] = useState<StatusFilter>("");
    const [priority, setPriority] = useState<PriorityFilter>("");

    const token = useMemo(() => getToken(), []);

    useEffect(() => {
        if (!token) {
            router.replace("/login");
            return;
        }

        const fetchTickets = async () => {
            try {
                setLoading(true);

                const params = new URLSearchParams();
                if (status) params.set("status", status);
                if (priority) params.set("priority", priority);
                params.set("page", "1");
                params.set("limit", "20");

                const response = await apiRequest<ListTicketsResponse>(
                    `/tickets?${params.toString()}`,
                    {
                        method: "GET",
                        token,
                    }
                );

                setTickets(response.data);
            } catch (error) {
                const err = error as Error;
                message.error(err.message || "Failed to fetch tickets");
            } finally {
                setLoading(false);
            }
        };

        fetchTickets();
    }, [token, router, status, priority]);

    const handleLogout = () => {
        removeToken();
        router.push("/login");
    };

    const columns: ColumnsType<Ticket> = [
        {
            title: "ID",
            dataIndex: "id",
            key: "id",
            width: 80,
        },
        {
            title: "Title",
            dataIndex: "title",
            key: "title",
            render: (value: string) => <Text strong>{value}</Text>,
        },
        {
            title: "Status",
            dataIndex: "status",
            key: "status",
            width: 140,
            render: (value: Ticket["status"]) => {
                switch (value) {
                    case "open":
                        return <Tag>Open</Tag>;
                    case "in_progress":
                        return <Tag color="gold">In Progress</Tag>;
                    case "resolved":
                        return <Tag color="green">Resolved</Tag>;
                    default:
                        return <Tag>{value}</Tag>;
                }
            },
        },
        {
            title: "Priority",
            dataIndex: "priority",
            key: "priority",
            width: 120,
            render: (value: Ticket["priority"]) => {
                switch (value) {
                    case "low":
                        return <Tag>Low</Tag>;
                    case "medium":
                        return <Tag color="orange">Medium</Tag>;
                    case "high":
                        return <Tag color="red">High</Tag>;
                    default:
                        return <Tag>{value}</Tag>;
                }
            },
        },
        {
            title: "Created By",
            dataIndex: "created_by",
            key: "created_by",
            width: 120,
        },
        {
            title: "Assigned To",
            dataIndex: "assigned_to",
            key: "assigned_to",
            width: 120,
            render: (value: number | null) => value ?? "-",
        },
        {
            title: "Updated At",
            dataIndex: "updated_at",
            key: "updated_at",
            width: 220,
            render: (value: string) => new Date(value).toLocaleString(),
        },
    ];

    return (
        <div style={{ padding: 24 }}>
            <Card>
                <Flex justify="space-between" align="center" style={{ marginBottom: 24 }}>
                    <div>
                        <Title level={2} style={{ margin: 0 }}>
                            Tickets
                        </Title>
                        <Text type="secondary">Support dashboard ticket list</Text>
                    </div>

                    <Button danger onClick={handleLogout}>
                        Logout
                    </Button>
                </Flex>

                <Space style={{ marginBottom: 16 }} wrap>
                    <Select
                        value={status}
                        onChange={(value) => setStatus(value)}
                        style={{ width: 180 }}
                        options={[
                            { value: "", label: "All Statuses" },
                            { value: "open", label: "Open" },
                            { value: "in_progress", label: "In Progress" },
                            { value: "resolved", label: "Resolved" },
                        ]}
                    />

                    <Select
                        value={priority}
                        onChange={(value) => setPriority(value)}
                        style={{ width: 180 }}
                        options={[
                            { value: "", label: "All Priorities" },
                            { value: "low", label: "Low" },
                            { value: "medium", label: "Medium" },
                            { value: "high", label: "High" },
                        ]}
                    />
                </Space>

                <Table<Ticket>
                    rowKey="id"
                    columns={columns}
                    dataSource={tickets}
                    loading={loading}
                    pagination={false}
                    onRow={(record) => ({
                        onClick: () => router.push(`/tickets/${record.id}`),
                        style: { cursor: "pointer" },
                    })}
                />
            </Card>
        </div>
    );
}