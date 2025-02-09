import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table"
import {ContainerStatus} from "@/types/container-status.ts";
import {Card, CardContent, CardHeader, CardTitle} from "./ui/card";


interface ContainerTableProps {
    containerStatuses: ContainerStatus[]
}

export function ContainerTable({containerStatuses}: ContainerTableProps) {
    return (
        <div className="max-w-4xl mx-auto p-6 w-screen h-screen">
            <Card className="">
                <CardHeader className="space-y-4">
                    <div className="flex items-center justify-center gap-3">
                        <div>
                            <CardTitle className="text-2xl font-bold text-blue-900">Docker Container Monitor</CardTitle>
                            <p className="text-sm text-blue-600">Real-time container status</p>
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="rounded-lg overflow-hidden border bg-white">
                        <Table>
                            <TableHeader>
                                <TableRow >
                                    <TableHead className="w-[200px] font-semibold text-center">IP Address</TableHead>
                                    <TableHead className="font-semibold text-center">Ping Time (ms)</TableHead>
                                    <TableHead className="text-right font-semibold">Last Ping</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {containerStatuses.map((status, index) => (
                                    <TableRow
                                        key={index}
                                        className="hover:bg-blue-50/50 transition-colors"
                                    >
                                        <TableCell className="font-medium text-center">{status.ip}</TableCell>
                                        <TableCell className="text-center">{+status.ping_time.toFixed(6)}</TableCell>
                                        <TableCell className="text-right text-gray-600">
                                            {new Date(status.last_ping).toLocaleString()}
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}