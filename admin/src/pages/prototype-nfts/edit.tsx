import React, { useEffect } from "react";
import { IResourceComponentsProps, useApiUrl } from "@refinedev/core";
import { Edit, useForm, getValueFromEvent } from "@refinedev/antd";
import { Button, Form, Input, Space, Upload } from "antd";

export const PrototypeNftsEdit: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps, onFinish, queryResult } = useForm();

  const [attributes, setAttributes] = React.useState<
  {
    trait_type: string;
    value: string;
  }[]
>([]);

  const handleOnFinish = (values: any) => {

    onFinish({
      ...values,
      weight: Number(values.weight),
      attributes: attributes,
    });
  };


  useEffect(() => {
    if (attributes.length === 0) {
    setAttributes(queryResult?.data?.data?.attributes || []);
    }
  }, [queryResult?.isSuccess]);

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} onFinish={handleOnFinish} layout="vertical">
        <Form.Item
          label="Id"
          name={["id"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
        <Form.Item
          label="Base64"
          name={["base64"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
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
          label="External Url"
          name={["external_url"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item label="Image">
          <Form.Item
            name="image"
            getValueProps={(value) => ({
              fileList: [{ url: value, name: value, uid: value }],
            })}
            getValueFromEvent={getValueFromEvent}
            noStyle
            rules={[
              {
                required: true,
              },
            ]}
          >
            <Upload.Dragger listType="picture" beforeUpload={() => false}>
              <p className="ant-upload-text">Drag & drop a file in this area</p>
            </Upload.Dragger>
          </Form.Item>
        </Form.Item>
        <Form.Item label="Weight" name={["weight"]}>
          <Input />
        </Form.Item>

        <Form.Item
          label="Marketplace"
          name={["marketplace"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <h2>Attributes</h2>
        {attributes.map((attr, index) => (
          <div key={index}>
            <Form.Item label={`Key ${index + 1}`}>
              <Input
                size="large"
                placeholder={`Key ${index + 1}`}
                value={attr.trait_type}
                onChange={(e) =>
                  setAttributes(
                    attributes.map((item, i) =>
                      i === index
                        ? { ...item, trait_type: e.target.value }
                        : item
                    )
                  )
                }
              />
            </Form.Item>
            <Form.Item label={`Value ${index + 1}`}>
              <Input
                size="large"
                placeholder={`Value ${index + 1}`}
                value={attr.value}
                onChange={(e) =>
                  setAttributes(
                    attributes.map((item, i) =>
                      i === index ? { ...item, value: e.target.value } : item
                    )
                  )
                }
              />
            </Form.Item>
          </div>
        ))}
        <Space size="large">
          <Button
            type="primary"
            onClick={() => {
              // add empty attribute
              setAttributes([
                ...attributes,
                {
                  trait_type: "",
                  value: "",
                },
              ]);
            }}
          >
            Add Attribute
          </Button>
          {attributes.length > 0 && (
            <Button
              onClick={() =>
                // remove last attribute
                setAttributes(attributes.slice(0, attributes.length - 1))
              }
            >
              Remove Attribute
            </Button>
          )}
        </Space>
      </Form>
    </Edit>
  );
};
