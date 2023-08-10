import React from "react";
import { IResourceComponentsProps } from "@refinedev/core";
import { Edit, useForm, getValueFromEvent } from "@refinedev/antd";
import { Form, Image, Input, Upload } from "antd";

export const CollectionsEdit: React.FC<IResourceComponentsProps> = () => {
    const { formProps, saveButtonProps, queryResult } = useForm();

    const collectionsData = queryResult?.data?.data;

    return (
        <Edit saveButtonProps={saveButtonProps}>
            <Form {...formProps} layout="vertical">
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
                    label="Raw Address"
                    name={["raw_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Friendly Address"
                    name={["friendly_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Next Item Index"
                    name={["next_item_index"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Content Uri"
                    name={["content_uri"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Raw Owner Address"
                    
                    name={["raw_owner_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Friendly Owner Address"
                    name={["friendly_owner_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input disabled readOnly />
                </Form.Item>
                <Form.Item
                    label="Name"
                    name={["name"]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Description"
                    name={["description"]}
                >
                    <Input />
                </Form.Item>
                <Form.Item label="Image">
                    <Form.Item
                        name="image"
                        getValueFromEvent={getValueFromEvent}
                        noStyle
                    >
                        <Input />
                    </Form.Item>
                    <div>
                    <Image src={collectionsData?.image} />
                    </div>
                </Form.Item>
                <Form.Item
                    label="Default rating points"
                    name={["default_weight"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Created At"
                    name={["created_at"]}
                >
                    <Input disabled />
                </Form.Item>
                <Form.Item
                    label="Updated At"
                    name={["updated_at"]}
                >
                    <Input disabled />
                </Form.Item>
                <Form.Item
                    label="Version"
                    name={["version"]}
                >
                    <Input disabled />
                </Form.Item>
            </Form>
        </Edit>
    );
};