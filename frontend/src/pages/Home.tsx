import { useEffect, useState } from "react"
import { ContainerTable } from "@/components/container-table"
import {ContainerStatus} from "@/types/container-status.ts";
import {useToast} from "@/hooks/use-toast.ts";

export default function Home() {
    const [containerStatuses, setContainerStatuses] = useState<ContainerStatus[]>([])
    const { toast } = useToast()

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch("/api/containers")
                const data = await response.json()
                setContainerStatuses(data)
            } catch (error) {
                console.error("Error fetching data:", error)

                toast({
                    variant: "destructive",
                    title: "Error fetching data.",
                    description: "There was a problem during your request.",
                })
            }
        }

        fetchData()
        const interval = setInterval(fetchData, 10000)

        return () => clearInterval(interval)
    }, [])

    return (
        <div className="container mx-auto py-10">
            <ContainerTable containerStatuses={containerStatuses} />
        </div>
    )
}