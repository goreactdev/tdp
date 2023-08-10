import React from "react";
import { IResourceComponentsProps, useList } from "@refinedev/core";
import { Create, useForm } from "@refinedev/antd";
import { Form, Input, Checkbox, Space, Typography, AutoComplete } from "antd";
import { useState, useEffect } from "react";
import { SearchOutlined } from "@ant-design/icons";

const { Text } = Typography;

const renderTitle = (title: string) => {
  return (
    <Text strong style={{ fontSize: "16px" }}>
      {title}
    </Text>
  );
};

// To be able to customize the option item
const renderItem = (firstText: string, secondText: string) => {
  return {
    value: firstText,
    label: (
      <Space.Compact direction="vertical">
        <Text strong>{firstText}</Text>
        <Text>{secondText}</Text>
      </Space.Compact>
    ),
  };
};

export const ActivitiesCreate: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps, onFinish } = useForm<{
    name: string;
    description: string;
    token_threshold: number;
  }>();

  const [searchData, setSearchData] = useState<{
    token_address: string;
  }>({
    token_address: "",
  });

  const [tokenOptions, setTokenOptions] = useState<any>([]);

  const { refetch: refetchTokens } = useList({
    resource: "prototype-nfts",
    filters: [
      {
        field: "name",
        operator: "contains",
        value: searchData.token_address,
      },
    ],
    queryOptions: {
      enabled: false,
      onSuccess: (data) => {
        const tokensOptionGroup = data.data.map((item) =>
          renderItem(item.base64, item.display_name)
        );
        if (tokensOptionGroup.length > 0) {
          setTokenOptions([
            {
              label: renderTitle("Tokens"),
              options: tokensOptionGroup,
            },
          ]);
        }
      },
    },
  });

  useEffect(() => {
    refetchTokens();
    setTokenOptions([]);
  }, [searchData.token_address]);

  const handleOnFinish = (values: any) => {
    onFinish( {
      ...values,
      token_threshold: Number(values.token_threshold),
    })
  }

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} onFinish={handleOnFinish} layout="vertical">
        <Form.Item
          label="SBT Token Metadata"
          name={["sbt_token_metadata"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <AutoComplete
            style={{ width: "100%", maxWidth: "550px" }}
            filterOption={false}
            options={tokenOptions}
            value={searchData.token_address}
            onChange={(value: string) =>
              setSearchData({ ...searchData, token_address: value })
            }
            onSearch={(value: string) =>
              setSearchData({ ...searchData, token_address: value })
            }
          >
            <Input
              size="large"
              value={searchData.token_address}
              placeholder="Search activities"
              suffix={<SearchOutlined />}
            />
          </AutoComplete>
        </Form.Item>

        <Form.Item
          label="Name"
          name={["name"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Description"
          name={["description"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Token Threshold"
          name={["token_threshold"]}
          rules={[
            {
              required: true,
            },
            {
              min: 2
            },
          ]}
        >
          <Input />
        </Form.Item>
      </Form>
    </Create>
  );
};
