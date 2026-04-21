// components/RecruiterHint.tsx
import { Typography, Space, Tag } from "antd";
import { InfoCircleOutlined } from "@ant-design/icons";

const { Text } = Typography;

type RecruiterHintProps = {
  text: string;
  tags?: string[];
};

export default function RecruiterHint({
  text,
  tags = [],
}: RecruiterHintProps) {
  return (
    <div
      style={{
        marginTop: 10,
        padding: "10px 12px",
        borderLeft: "3px solid #1677ff",
        background: "rgba(22, 119, 255, 0.04)",
        borderRadius: 8,
      }}
    >
      <Space orientation="vertical" size={6} style={{ width: "100%" }}>
        <Space size={6}>
          <InfoCircleOutlined style={{ color: "#1677ff" }} />
          <Text strong style={{ fontSize: 13 }}>
            Recruiter note
          </Text>
        </Space>

        <Text style={{ 
            fontSize: 13, 
            color: "rgba(0,0,0,0.72)",
            whiteSpace: "pre-line", 
        }}>{text}</Text>

        {tags.length > 0 && (
          <Space wrap size={[6, 6]}>
            {tags.map((tag) => (
              <Tag key={tag}>{tag}</Tag>
            ))}
          </Space>
        )}
      </Space>
    </div>
  );
}