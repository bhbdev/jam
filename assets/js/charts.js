/**
 * For usage, visit Chart.js docs https://www.chartjs.org/docs/latest/
 */

function NewPieChart(container, data, labels, colors) {
    const pieConfig = {
        type: 'doughnut',
        data: {
            datasets: [
                {
                    data: data,
                    backgroundColor: colors,
                    label: 'Applications',
                },
            ],
            labels: labels,
        },
        options: {
            responsive: true,
            cutoutPercentage: 80,
            legend: {
                display: false,
            },
        },
    }
    const pieCtx = document.getElementById(container)
    window.myPie = new Chart(pieCtx, pieConfig)
}

function NewBarChart(container, data, labels, colors) {
    const barConfig = {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Dataset 1',
                    backgroundColor: colors,
                    data: data,
                },
            ],
        },
        options: {
            responsive: true,
            legend: {
                display: false,
            },
        },
    }
    const barCtx = document.getElementById(container)
    window.myBar = new Chart(barCtx, barConfig)
}

function NewLineChart(container, data, labels, colors) {
    const lineConfig = {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Dataset 1',
                    backgroundColor: colors,
                    borderColor: colors,
                    data: data,
                    fill: false,
                },
            ],
        },
        options: {
            responsive: true,
            legend: {
                display: false,
            },
            tooltips: {
                mode: 'index',
                intersect: false,
            },
            hover: {
                mode: 'nearest',
                intersect: true,
            },
            scales: {
                x: {
                  display: true,
                  scaleLabel: {
                    display: true,
                    labelString: 'Month',
                  },
                },
                y: {
                  display: true,
                  scaleLabel: {
                    display: true,
                    labelString: 'Value',
                  },
                },
            },
        },
    }
    const lineCtx = document.getElementById(container)
    window.myLine = new Chart(lineCtx, lineConfig)
}